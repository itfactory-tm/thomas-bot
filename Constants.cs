using DSharpPlus.Entities;
using System;

namespace ThomasBot
{
    class Constants
    {
        //public const ulong botChannel = 0; // id of the bot channel

        public const string rolename_moderator = "Moderator";
        public const string rolename_admin = "Admin";

        public static DiscordColor GetColor(ConstColors color)
        {
            return color switch
            {
                ConstColors.CommandRun => DiscordColor.VeryDarkGray,
                ConstColors.Error => DiscordColor.Red,
                ConstColors.Warning => DiscordColor.Yellow,
                ConstColors.LogGeneric => DiscordColor.VeryDarkGray,
                ConstColors.LogWarning => DiscordColor.Yellow,
                ConstColors.LogDangerous => DiscordColor.Red,
                _ => throw new NotImplementedException("This color hasn't been defined.")
            };
        }

        public enum ConstColors
        {
            CommandRun,
            Error,
            Warning,
            LogGeneric,
            LogWarning,
            LogDangerous
        }
    }
}
