using DSharpPlus;
using DSharpPlus.CommandsNext;
using DSharpPlus.Entities;
using DSharpPlus.Exceptions;
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using ThomasBot.Exceptions;

namespace ThomasBot
{
    /// <summary>
    /// Class providing various extension methods.
    /// </summary>
    static class Extensions
    {
        /// <summary>
        /// Returns the full invite url of an invite. e.g. https://discord.gg/2jk86sa8
        /// </summary>
        /// <param name="invite">The invite object to get the invite code from.</param>
        /// <returns></returns>
        public static string GetFullUrl(this DiscordInvite invite)
        {
            return $"https://discord.gg/{invite.Code} "; //that space at the end is very important!
        }

        /// <summary>
        /// Converts a hexadecimal string to an integer.
        /// </summary>
        /// <param name="x">The hexadecimal string to convert. Allowed characters are 0-9 a-f and A-F.</param>
        /// <returns>The integer representation of the given hex value.</returns>
        public static int HexToInt(this string x)
        {
            char[] allowedChars = { '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f' };
            int retval = 0;
            x = x.ToLower();
            if (x.Any(y => !allowedChars.Contains(char.ToLower(y))))
            {
                throw new ArgumentException();
            }
            for (int i = 1; i <= x.Length; i++)
            {
                retval += (int)Math.Pow(16, i - 1) * allowedChars.ToList().IndexOf(x[^i]);
            }
            return retval;
        }

        /// <summary>
        /// Grabs a random element of a generic IEnumberalbe.
        /// </summary>
        /// <typeparam name="T">Generic type used for the operation.</typeparam>
        /// <param name="enumerable">The object storing the elements.</param>
        /// <returns>A random element thats contained in the input.</returns>
        public static T Random<T>(this IEnumerable<T> enumerable)
        {
            var list = enumerable as IList<T> ?? enumerable.ToList();
            return list.ElementAt(ThomasBot.rand.Next(0, list.Count()));
        }

        /// <summary>
        /// Grabs a random element of a generic IEnumberalbe.
        /// </summary>
        /// <typeparam name="T">Generic type used for the operation.</typeparam>
        /// <param name="enumerable">The object storing the elements.</param>
        /// <param name="weight">A function returning a weight for every elemment. The bigger the more likely it is that this element is chosen.</param>
        /// <returns>A random element thats contained in the input.</returns>
        public static T WeightedRandom<T>(this IEnumerable<T> enumerable, Func<T, int> weight)
        {
            List<T> list = new List<T>();
            foreach (T t in enumerable)
            {
                for (int i = 0; i < weight(t); i++)
                {
                    list.Add(t);
                }
            }

            return list.ElementAt(ThomasBot.rand.Next(0, list.Count()));
        }

        /// <summary>
        /// Checks if this message mentioned the bot.
        /// </summary>
        /// <param name="msg">The message to check.</param>
        /// <returns>A boolean describing whether the bot was mentioned or not.</returns>
        public static bool MentionedMe(this DiscordMessage msg)
        {
            try
            {
                if (msg.MentionedUsers.Count(x => x.Id == msg.Channel.Guild.CurrentMember.Id) > 0)
                {
                    return true;
                }
                //if (msg.MentionedRoles.Intersect(msg.Channel.Guild.Members.First(x => x.Id == 327150443610505227).Roles).Count() > 0)
                //{
                //    return true;
                //}
                return false;
            }
            catch (Exception ex)
            {
                return false;
                throw ex;
            }
        }

        /// <summary>
        /// Method for performing multiple replaces at once and returning the new String.
        /// </summary>
        /// <typeparam name="T">Generic type used.</typeparam>
        /// <param name="str">The IEnumerable that will be searched.</param>
        /// <param name="search">Item to look for.</param>
        /// <param name="newobj">Item to replace the found one with. (corresponding to the 'search' argument.</param>
        /// <returns>An IEnumberable where all the matched items were replaced.</returns>
        public static IEnumerable<T> ReplaceMulti<T>(this IEnumerable<T> str, T[] search, T newobj)
        {
            foreach (var item in str)
            {
                yield return search.Contains(item) ? newobj : item;
            }
        }

        public static string ReplaceMulti(this string str, string[] search, string newstr)
        {
            return search.Contains(str) ? newstr : str;
        }

        public static string ToReadableString(this TimeSpan span, bool superCompact = false, bool replaceLastComma = false)
        {
            string formatted;

            if (!superCompact)
            {
                formatted = string.Format("{0}{1}{2}{3}{4}",
                        span.Duration().Days >= 365 ? string.Format("{0:0} year{1}, ", span.Days / 365, span.Days / 365 == 1 ? string.Empty : "s") : string.Empty,
                        span.Duration().Days % 365 > 0 ? string.Format("{0:0} day{1}, ", span.Days % 365, span.Days % 365 == 1 ? string.Empty : "s") : string.Empty,
                        span.Duration().Hours > 0 ? string.Format("{0:0} hour{1}, ", span.Hours, span.Hours == 1 ? string.Empty : "s") : string.Empty,
                        span.Duration().Minutes > 0 ? string.Format("{0:0} minute{1}, ", span.Minutes, span.Minutes == 1 ? string.Empty : "s") : string.Empty,
                        span.Duration().Seconds > 0 ? string.Format("{0:0} second{1}", span.Seconds, span.Seconds == 1 ? string.Empty : "s") : string.Empty);

                if (formatted.EndsWith(", ")) formatted = formatted[0..^2];

                if (string.IsNullOrEmpty(formatted)) formatted = "0 seconds";

                if (replaceLastComma && formatted.Split(',').Count() > 1)
                {
                    formatted = (string.Join(",", formatted.Split(',').SkipLast(1)) + $" and {formatted.Split(',').Last()}").Replace("  ", " ");
                }
            }
            else
            {
                formatted = string.Format("{0}{1}{2}{3}{4}{5}",
                        span.Duration().Days >= 365 ? string.Format("{0:0}y, ", span.Days / 365) : string.Empty,
                        span.Duration().Days % 365 / (365 / 12f) >= 1 ? string.Format("{0:0}M, ", Math.Floor(span.Days % 365 / (365 / 12f))) : string.Empty,
                        span.Duration().Days % 365 % (365 / 12f) >= 1 ? string.Format("{0:0}d, ", span.Days % 365 % (365 / 12f)) : string.Empty,
                        span.Duration().Hours > 0 ? string.Format("{0:0}h, ", span.Hours) : string.Empty,
                        span.Duration().Minutes > 0 ? string.Format("{0:0}m, ", span.Minutes) : string.Empty,
                        span.Duration().Seconds > 0 ? string.Format("{0:0}s", span.Seconds) : string.Empty);

                if (formatted.EndsWith(", ")) formatted = formatted[0..^2];

                if (string.IsNullOrEmpty(formatted)) formatted = "0s";
            }
            return formatted;
        }

        public static async Task<DiscordMember> ToMemberAsync(this DiscordUser user, DiscordGuild guild, bool doThrow = true)
        {
            if (guild == null)
            {
                if (doThrow)
                    throw new MemberNotFoundException("Could not convert user to member because the specified guild was null.");
                else
                    return null;
            }

            try
            {
                var result = guild.Members.FirstOrDefault(x => x.Key == user.Id).Value ?? (await guild.GetMemberAsync(user.Id));

                if (result == null && doThrow)
                {
                    throw new MemberNotFoundException("Could not convert user to member because the specified member does not exist.");
                }
                return result;
            }
            catch (NotFoundException)
            {
                if (doThrow)
                    throw new MemberNotFoundException("Could not convert user to member because the specified member does not exist in this guild.");

                return null;
            }

        }

        public static DiscordMember ToMemberNonAsync(this DiscordUser user, DiscordGuild guild, bool doThrow = true)
        {
            if (guild == null)
            {
                if (doThrow)
                    throw new MemberNotFoundException("Could not convert user to member because the specified guild was null.");
                else
                    return null;
            }

            var result = guild.Members.FirstOrDefault(x => x.Key == user.Id).Value ?? (guild.GetMemberAsync(user.Id).GetAwaiter().GetResult());
            if (result == null && doThrow)
            {
                throw new MemberNotFoundException("Could not convert user to member because the specified member does not exist.");
            }
            return result;
        }

        public static bool IsModOrAdmin(this DiscordMember member)
        {
            return member.Roles.Any(x => x.Name == Constants.rolename_moderator || x.CheckPermission(Permissions.Administrator) == PermissionLevel.Allowed);
        }
        public static bool IsAdmin(this DiscordMember member)
        {
            return member.Roles.Any(x => x.CheckPermission(Permissions.Administrator) == PermissionLevel.Allowed);
        }

        /// <summary>
        /// Checks if the given user is able to perform actions that require the given permission. Also takes Administrator into account.
        /// </summary>
        /// <param name="member"></param>
        /// <param name="perm"></param>
        /// <param name="allowCreatorOverride"></param>
        /// <returns></returns>
        public static bool HasPermission(this DiscordMember member, Permissions perm)
        {
            return member.Roles.Any(x => x.CheckPermission(perm) == PermissionLevel.Allowed) || member.IsAdmin();
        }
        public static string GetUsernameNickAndID(this DiscordMember member)
        {
            return $"{member.Username} {((member?.Nickname != null) ? "[" + member.Nickname + "]" : "")} [{member.Id}]";
        }
        public static IEnumerable<T> SkipLast<T>(this IEnumerable<T> e, int amount)
        {
            return e.Take(e.Count() - amount);
        }

        public static List<string> ReadAllLines(this FileStream fs)
        {
            var retval = new List<string>();
            using (var sr = new StreamReader(fs, Encoding.Default, true, 1024, true))
            {
                string line;
                while ((line = sr.ReadLine()) != null)
                {
                    retval.Add(line);
                }
            }
            return retval;
        }

        public static void WriteAllLines(this FileStream fs, IEnumerable<string> lines)
        {
            using (var sw = new StreamWriter(fs, Encoding.Default, 1024, true))
            {
                sw.BaseStream.SetLength(0);
                sw.BaseStream.Seek(0, SeekOrigin.Begin);
                if (lines.Count() == 0)
                {
                    sw.WriteLine("");
                }
                else
                {
                    foreach (var line in lines)
                    {
                        sw.WriteLine(line);
                    }
                }
            }
        }

        public static void AppendAllLines(this FileStream fs, IEnumerable<string> lines)
        {
            var allLines = new List<string>();
            using (var sr = new StreamReader(fs, Encoding.Default, true, 1024, true))
            {
                string line;
                while ((line = sr.ReadLine()) != null)
                {
                    allLines.Add(line);
                }
            }
            allLines.AddRange(lines);
            using (var sw = new StreamWriter(fs, Encoding.Default, 1024, true))
            {
                sw.BaseStream.Seek(0, SeekOrigin.End);
                foreach (var line in allLines)
                {
                    sw.WriteLine(line);
                }
            }
        }

        public static IEnumerable<TSource> DistinctBy<TSource, TKey>(this IEnumerable<TSource> source, Func<TSource, TKey> keySelector)
        {
            HashSet<TKey> seenKeys = new HashSet<TKey>();
            foreach (TSource element in source)
            {
                if (seenKeys.Add(keySelector(element)))
                {
                    yield return element;
                }
            }
        }

        public static async Task<DiscordMessage> RespondWithEmbedAsync(this CommandContext ctx, string title, string text, string footer = null, DiscordColor? customColor = null, string footerIconUrl = null, DateTime? timestamp = null)
        {
            return await ctx.Channel.SendEmbedAsync(title, text, footer, customColor, footerIconUrl, timestamp);
        }

        public static async Task<DiscordMessage> SendEmbedAsync(this DiscordChannel chan, string title, string text, string footer = null, DiscordColor? customColor = null, string footerIconUrl = null, DateTime? timestamp = null)
        {
            DiscordColor color = Constants.GetColor(Constants.ConstColors.CommandRun);
            if (customColor != null)
            {
                color = customColor.Value;
            }
            return await chan.SendMessageAsync(embed: Factory.GetEmbed(color, title, text, footer, footerIconUrl, timestamp));
        }

        public static string ToUsernameAndDiscriminatorString(this DiscordUser user)
        {
            if (user != null)
            {
                return user.Username + "#" + user.Discriminator;
            }
            throw new ArgumentNullException("Argument 'user' was null!");
        }
    }
}
