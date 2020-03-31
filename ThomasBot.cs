using DSharpPlus;
using DSharpPlus.CommandsNext;
using DSharpPlus.CommandsNext.Attributes;
using DSharpPlus.CommandsNext.Exceptions;
using DSharpPlus.Entities;
using DSharpPlus.EventArgs;
using DSharpPlus.Interactivity;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using ThomasBot.Commands;
using ThomasBot.Exceptions;

namespace ThomasBot
{
    class ThomasBot
    {
        public const string BotVersion = "Thomasbot 1.0";

        public static bool ShutdownRequested { get; private set; }

        public static Random rand = new Random();
        public static DateTime startTime = DateTime.UtcNow;
        public static long commandsRanThisSess = 0;

        public DiscordClient Client { get; private set; }
        public string CommandPrefix { get; private set; }
        public bool commandHandlerEn = false;
        public CommandsNextExtension Commands { get; private set; }
        public InteractivityExtension Interactivity { get; private set; }

        internal async Task RunBotAsync()
        {
            Thread.CurrentThread.CurrentCulture = new System.Globalization.CultureInfo("en-US");
            AppDomain.CurrentDomain.UnhandledException += CurrentDomain_UnhandledException;

            string configFile = "config.json";
            if (!File.Exists(configFile))
            {
                Log.WriteLogMessage($"Unable to load find file '{configFile}'. It was now automatically created. Please fill in the contents.", LogOutputLevel.Critical);
                File.WriteAllText(configFile, JsonConvert.SerializeObject(new ConfigJson()), Encoding.Default);
                return;
            }

            var json = File.ReadAllText(configFile, Encoding.Default);
            var cfgjson = JsonConvert.DeserializeObject<ConfigJson>(json);

            LogLevel desiredLogLevel = LogLevel.Info;

#if DEBUG
            desiredLogLevel = LogLevel.Debug;
            Log.WriteLogMessage("--------------WARNING--------------", LogOutputLevel.Warning);
            Log.WriteLogMessage("-------RUNNING IN DEBUG MODE-------", LogOutputLevel.Warning);
            Log.WriteLogMessage("--------------WARNING--------------", LogOutputLevel.Warning);
#endif

            var cfg = new DiscordConfiguration
            {
                Token = cfgjson.Token,

                TokenType = TokenType.Bot,

                AutoReconnect = true,
                LogLevel = desiredLogLevel,
                UseInternalLogHandler = true,
            };

            this.Client = new DiscordClient(cfg);

            // next, let's hook some events, so we know
            // what's going on
            this.Client.Ready += Client_ReadyAsync;
            this.Client.ClientErrored += Client_ClientError;

            this.CommandPrefix = cfgjson.CommandPrefix;

            // up next, let's set up our commands
            var ccfg = new CommandsNextConfiguration
            {
                StringPrefixes = new[] { this.CommandPrefix },

                // enable responding in direct messages
                EnableDms = true,

                // enable mentioning the bot as a command prefix
                EnableMentionPrefix = true,

                EnableDefaultHelp = false,

                IgnoreExtraArguments = true
            };

            // and hook them up
            this.Commands = this.Client.UseCommandsNext(ccfg);
            this.Client.DebugLogger.LogMessageReceived += DebugLogger_LogMessageReceived;

            // let's hook some command events, so we know what's 
            // going on
            this.Commands.CommandExecuted += Commands_CommandExecuted;
            this.Commands.CommandErrored += Commands_CommandErrored;

            // up next, let's register our commands
            this.Commands.RegisterCommands<AdminCommands>();
            this.Commands.RegisterCommands<UserCommands>();


            var iConfig = new InteractivityConfiguration
            {

            };
            this.Interactivity = this.Client.UseInteractivity(iConfig);

            // Load database
            // Add DB if required
            // SaveLoad.LoadAll();
            // Log.WriteLogMessage("Database loaded.", LogOutputLevel.Info);

            // Clean database
            // SaveLoad.DoCleanup();
            Log.WriteLogMessage("Database cleaned.", LogOutputLevel.Info);

            // finally, let's connect and log in
            await this.Client.ConnectAsync();

            Factory.StartAsNewThread(async delegate
            {
                await ReadConsoleAsync();
            });

            // and this is to prevent premature quitting
            await Task.Delay(-1);
        }

        private Task ReadConsoleAsync()
        {
            string concmd = string.Empty;

            while (!this.commandHandlerEn)
            {
                Task.Delay(500);
            }

            Thread t = new Thread((ThreadStart)delegate { Task.Yield(); }); ;
            while (!ShutdownRequested)
            {
                concmd = Console.ReadLine();
                try
                {
                    bool cmdNotRun = false;
                    switch (concmd.ToLower())
                    {
                        case "exit": Shutdown(ShutdownAction.Shutdown); break;
                        case "reconnect": this.Client.ReconnectAsync(); break;
                        //case "save": SaveLoad.SaveAll(); break;
                        //case "load": SaveLoad.LoadAll(); break;
                        //case "save-on": Log.WriteLogMessage("Last state: enable Autosave: " + (CronStore.enableAutosave ? "TRUE" : "FALSE"), LogOutputLevel.Info); CronStore.enableAutosave = true; break;
                        //case "save-off": Log.WriteLogMessage("Last state: enable Autosave: " + (CronStore.enableAutosave ? "TRUE" : "FALSE"), LogOutputLevel.Info); CronStore.enableAutosave = false; break;

                        case "setdnd": this.Client.UpdateStatusAsync(userStatus: UserStatus.DoNotDisturb).GetAwaiter().GetResult(); break;
                        case "setonline": this.Client.UpdateStatusAsync(userStatus: UserStatus.Online).GetAwaiter().GetResult(); break;
                        case "setidle": this.Client.UpdateStatusAsync(userStatus: UserStatus.Idle).GetAwaiter().GetResult(); break;
                        case "setinvisible": this.Client.UpdateStatusAsync(userStatus: UserStatus.Invisible).GetAwaiter().GetResult(); break;

                        default: cmdNotRun = true; Console.WriteLine("Unknown command."); break;
                    }
                    if (!cmdNotRun)
                        Console.WriteLine("Command was executed.");

                }
                catch (Exception ex)
                {
                    Log.WriteLogMessage($"Error while running console command: {ex.ToString()}", LogOutputLevel.Error);
                }
            }
            return Task.CompletedTask;
        }

        private async Task Client_ReadyAsync(ReadyEventArgs e)
        {
            // let's log the fact that this event occured
            e.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", "Client is ready to process events.", DateTime.Now);

#if DEBUG
            await e.Client.UpdateStatusAsync(userStatus: UserStatus.DoNotDisturb);
#else
            await e.Client.UpdateStatusAsync(userStatus: UserStatus.Online);
#endif
            e.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", "Client has finished booting - console commands are enabled now. Type 'help' to see commands and 'exit' to exit.", DateTime.Now);
            this.commandHandlerEn = true;
        }
        private Task Client_ClientError(ClientErrorEventArgs e)
        {
            // let's log the details of the error that just 
            // occured in our client
            e.Client.DebugLogger.LogMessage(LogLevel.Error, "Thomas Bot", $"Exception occured: {e.Exception.GetType()}: {e.Exception.Message} {e.Exception.StackTrace}\n{e.Exception.InnerException} ({e.Exception.InnerException?.StackTrace})", DateTime.Now);

            // since this method is not async, let's return
            // a completed task, so that no additional work
            // is done
            return Task.CompletedTask;
        }

        private Task Commands_CommandExecuted(CommandExecutionEventArgs e)
        {
            // let's log the name of the command and user
            e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"{e.Context.User.Username} successfully executed '{e.Command.QualifiedName}'", DateTime.Now);

            return Task.CompletedTask;
        }

        private async Task Commands_CommandErrored(CommandErrorEventArgs e)
        {
#if !DEBUG
            if (e.Command?.QualifiedName == null) { return; }   //skip those command not found - messages
#endif

            // let's log the error details
            bool isError = true; // false = info; true = error
            // let's check if the error is a result of lack
            // of required permissions
            if (e.Exception is ChecksFailedException ex)
            {
                DiscordEmbedBuilder embed;
                IEnumerable<DSharpPlus.CommandsNext.Attributes.CooldownAttribute> cooldownAttribute;
                if (ex.FailedChecks.Any(x => x.GetType() == typeof(DSharpPlus.CommandsNext.Attributes.RequirePermissionsAttribute)))
                {
                    // yes, the user lacks required permissions, 
                    // let them know

                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");

                    // let's wrap the response into an embed
                    embed = new DiscordEmbedBuilder
                    {
                        Title = "Access denied",
                        Description = $"{emoji} You do not have the permissions required to execute this command.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"Error type: no permission", DateTime.Now);
                    isError = false; // show as info
                }
                // when an user tries to execute a staff command
                else if (ex.FailedChecks.Any(x => x.GetType() == typeof(RequireBotPermissionsAttribute)))
                {
                    // yes, the user lacks required permissions, 
                    // let them know

                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":warning:");

                    // let's wrap the response into an embed
                    embed = new DiscordEmbedBuilder
                    {
                        Title = "Missing permissions",
                        Description = $"{emoji} I do not have the permissions required to execute this command.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"Error type: no permission", DateTime.Now);
                    isError = false; // show as info
                }

                // if the command is only available to creators (eval)
                else if (ex.FailedChecks.Any(x => x.GetType() == typeof(DSharpPlus.CommandsNext.Attributes.RequireOwnerAttribute)))
                {
                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");

                    // let's wrap the response into an embed
                    embed = new DiscordEmbedBuilder
                    {
                        Title = "Access denied",
                        Description = $"{emoji} This command either is very dangerous or of no use to you. Thats why you aren't allowed to use it. Sorry!",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"Error type: no permission", DateTime.Now);
                    isError = false; // show as info
                }
                else if (ex.FailedChecks.Any(x => x.GetType() == typeof(Attributes.GuildOnly)))
                {
                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");

                    // let's wrap the response into an embed
                    embed = new DiscordEmbedBuilder
                    {
                        Title = "Not available",
                        Description = $"{emoji} This command is only available on servers and not in DMs. Join a server and send it in there.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"Error type: guild only - not in a guild", DateTime.Now);
                    isError = false; // show as info
                }
                else if ((cooldownAttribute = e.Command.ExecutionChecks?.Where(x => x.GetType() == typeof(DSharpPlus.CommandsNext.Attributes.CooldownAttribute))?.Select(x => ((DSharpPlus.CommandsNext.Attributes.CooldownAttribute)x)))?.Select(x => x.GetBucket(e.Context))?.FirstOrDefault(x => x.RemainingUses == 0) != null)
                {
                    var bucket = cooldownAttribute.First().GetBucket(e.Context);

                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");
                    embed = new DiscordEmbedBuilder
                    {
                        Title = "Slow down",
                        Description = $"{emoji} You reached the current usage limit of this command. You will be able to use it again in {bucket.Reset.ToReadableString()}.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"Error type: usages depleted - on cooldown. Cooldown Type: (per-){cooldownAttribute.First().BucketType.ToString()}", DateTime.Now);
                    isError = false; // show as info
                }
                else
                {
                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":question:");

                    // let's wrap the response into an embed
                    embed = new DiscordEmbedBuilder
                    {
                        Title = "Unknown Error",
                        Description = $"{emoji} An unknown error occurred.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    e.Context.Client.DebugLogger.LogMessage(LogLevel.Info, "Thomas Bot", $"Error type: unknown", DateTime.Now);
                }
                await e.Context.RespondAsync("", embed: embed);
            }

            else if (e.Exception is MemberNotFoundException)
            {
                var emoji = DiscordEmoji.FromName(e.Context.Client, ":grey_question:");

                // let's wrap the response into an embed
                var embed = new DiscordEmbedBuilder
                {
                    Title = "Error",
                    Description = $"{emoji} I was unable to find that user.",
                    Color = new DiscordColor(0xFF0000) // red
                };
                await e.Context.RespondAsync("", embed: embed);
                isError = false; // show as info
            }
            //TODO this may be an issue
            //check for if the casting failed. 
            else if (e.Exception is ArgumentException)
            {
                if (e.Command.Name == "eval")
                {

                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");
                    var embed = new DiscordEmbedBuilder
                    {
                        Title = "Error",
                        Description = $"Error: {e.Exception.ToString()}",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    await e.Context.RespondAsync("", embed: embed);
                }
                else if (e.Exception.Message != "Could not find a suitable overload for the command." && e.Command.Parent != null &&
                    e.Command.Overloads.Any(a => a.Arguments.Any(x => x.Type == typeof(DiscordMember) || x.Type == typeof(DiscordUser))))
                {
                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");

                    // let's wrap the response into an embed
                    var embed = new DiscordEmbedBuilder
                    {
                        Title = "Error",
                        Description = $"{emoji} That user can't be found.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    await e.Context.RespondAsync("", embed: embed);
                    isError = false; // show as info
                }
                else
                {
                    var emoji = DiscordEmoji.FromName(e.Context.Client, ":no_entry:");

                    // let's wrap the response into an embed
                    var embed = new DiscordEmbedBuilder
                    {
                        Title = "Error",
                        Description = $"{emoji} {e.Exception.Message} Type **>help {e.Command.Parent?.Name?.ToString() + " " + e.Command.Name}** to see the commands usage.",
                        Color = new DiscordColor(0xFF0000) // red
                    };
                    await e.Context.RespondAsync("", embed: embed);
                }
            }
            else if (e.Exception is DSharpPlus.Exceptions.UnauthorizedException)
            {

            }
            string debugInfo = "";
            if (isError)
            {
                debugInfo = $"\nCustom Stacktrace:'{e.Exception.StackTrace}'";
                if (e.Exception is DSharpPlus.Exceptions.BadRequestException)
                {
                    debugInfo += $"\nMessage:'{(e.Exception as DSharpPlus.Exceptions.BadRequestException).JsonMessage}'";
                }
                else
                {
                    debugInfo += $"\nMessage:'{e.Exception.Message}'";
                }
            }
            // show as info/error
            e.Context.Client.DebugLogger.LogMessage(isError ? LogLevel.Error : LogLevel.Info, "Thomas Bot", $"{e.Context.User.Username} tried executing '{e.Command?.QualifiedName ?? "<unknown command>"}' but it errored: {e.Exception.GetType()}: {e.Exception.Message ?? "<no message>"}" + $"\nCommand text: '{e.Context.Message.Content}'{debugInfo}", DateTime.Now);


        }

        private static bool shutdownAlreadyInitiated = false;
        public static void Shutdown(ShutdownAction action)
        {
            if (!shutdownAlreadyInitiated)
            {
                shutdownAlreadyInitiated = true;
                ShutdownRequested = true;
                Task.WaitAll(
                    //Task.Run(delegate { SaveLoad.SaveAll(); })
                    );
                Log.WriteLogMessage("Shutdown complete", LogOutputLevel.Info);
                Thread.Sleep(1000);

                switch (action)
                {
                    case ShutdownAction.Shutdown:
                        Environment.Exit(69);
                        break;
                    case ShutdownAction.Restart:
                        Environment.Exit(0);
                        break;
                    case ShutdownAction.Update:
                        Environment.Exit(70);
                        break;
                }

                Environment.Exit(69);
            }
        }

        public enum ShutdownAction
        {
            Shutdown,
            Restart,
            Update
        }
        private void CurrentDomain_UnhandledException(object sender, UnhandledExceptionEventArgs e)
        {
            string text = $"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] Unhandled Exception:\n" + e.ExceptionObject.ToString();
            File.AppendAllText("unhandled.txt", text + "\n");
        }

        // -- workaround for discord disconnecting us improperly --

        private int sessionStartAttemptFullRestartCounter = 4; // 4 = 4*30s = 2 minutes
        private void DebugLogger_LogMessageReceived(object sender, DebugLogMessageEventArgs e)
        {
            if (e.Level == LogLevel.Warning && e.Message == "Session start attempt was made while another session is active") // if we are in the session start attempt loop
            {
                this.sessionStartAttemptFullRestartCounter--;

                if (this.sessionStartAttemptFullRestartCounter <= 0)
                {
                    Log.WriteLogMessage("Session start attempt shenanigans are going on. Restarting!", LogOutputLevel.Warning);
                    Shutdown(ShutdownAction.Restart);
                }
            }
        }

        // --------------------------------------------------------

        public struct ConfigJson
        {
            [JsonProperty("token")]
            public string Token { get; private set; }

            [JsonProperty("prefix")]
            public string CommandPrefix { get; private set; }
        }
    }
}
