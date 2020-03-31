using DSharpPlus.Entities;
using System;
using System.Threading;

namespace ThomasBot
{
    class Factory
    {
        public static void DelayAction(int millisecond, Action action)
        {
            var timer = new Timer(delegate { Thread.Sleep(0); }, null, millisecond, Timeout.Infinite);
            timer = new Timer(delegate { action.Invoke(); timer.Dispose(); }, null, millisecond, Timeout.Infinite);
        }

        public static Thread StartAsNewThread(Action func)
        {
            Thread t = new Thread((ThreadStart)delegate { func.Invoke(); });
            t.Start();
            return t;
        }

        public static DiscordEmbed GetEmbed(DiscordColor color, string title, string text, string footer, string footerIconUrl = null, DateTime? timestamp = null)
        {
            var Builder = new DiscordEmbedBuilder();
            var b = Builder.WithColor(color).WithTitle(title).WithDescription(text).WithFooter(footer, footerIconUrl);
            if (timestamp.HasValue)
            {
                b.WithTimestamp(timestamp);
            }
            return b.Build();
        }

        public static DiscordEmbed GetRequestedByEmbed(DiscordColor color, string title, string text, DiscordUser user)
        {
            var Builder = new DiscordEmbedBuilder();
            return Builder.WithColor(color).WithTitle(title).WithDescription(text).WithFooter($"Requested by {user.Username}#{user.Discriminator}.", user.AvatarUrl);
        }
    }
}
