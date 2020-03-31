using DSharpPlus;
using DSharpPlus.CommandsNext;
using DSharpPlus.CommandsNext.Attributes;
using System;
using System.Linq;
using System.Threading.Tasks;

namespace ThomasBot
{
    class Attributes
    {
        /// <summary>
        /// Defines that usage of this command is restricted to Guilds.
        /// </summary>
        [AttributeUsage(AttributeTargets.Method | AttributeTargets.Class, AllowMultiple = false, Inherited = true)]
        public sealed class GuildOnly : CheckBaseAttribute
        {
            public override Task<bool> ExecuteCheckAsync(CommandContext ctx, bool help)
                => Task.FromResult(ctx.Guild != null || help);
        }

        /// <summary>
        /// Defines that usage of this command is restricted to the bot-channel. Developers can override this.
        /// </summary>
        [AttributeUsage(AttributeTargets.Method | AttributeTargets.Class, AllowMultiple = false, Inherited = true)]
        public sealed class BotChannelOnly : CheckBaseAttribute
        {
            public override Task<bool> ExecuteCheckAsync(CommandContext ctx, bool help)
                => Task.FromResult(ctx.Channel.IsPrivate /*|| ctx.Channel.Id == Constants.botChannel*/  || ctx.Member.IsAdmin() || ctx.Channel.PermissionOverwrites.Any(x => x.Type == OverwriteType.Role && x.Id == ctx.Guild.EveryoneRole.Id && x.Denied.HasPermission(Permissions.AccessChannels)));
        }  
    }
}
