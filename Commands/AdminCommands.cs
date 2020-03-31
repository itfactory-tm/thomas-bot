using DSharpPlus;
using DSharpPlus.CommandsNext;
using DSharpPlus.CommandsNext.Attributes;
using DSharpPlus.Entities;
using DSharpPlus.Interactivity;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using System.Reflection;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using static ThomasBot.Attributes;

namespace ThomasBot.Commands
{

    [RequirePermissions(Permissions.Administrator)]
    class AdminCommands : BaseCommandModule
    {
        [Command("mute")]
        [Description("Een gebruiker muten (admin only)")]
        [RequireBotPermissions(Permissions.ManageRoles)]
        public async Task Mute(CommandContext ctx, DiscordMember member)
        {
            var role = ctx.Guild.Roles.First(x => x.Value.Name == "Muted").Value;
            await member.GrantRoleAsync(role);
            await ctx.RespondAsync(":mute:");
        }

        [Command("unmute")]
        [Description("Een gebruiker unmuten (admin only)")]
        [RequireBotPermissions(Permissions.ManageRoles)]
        public async Task Unmute(CommandContext ctx, DiscordMember member)
        {
            var role = ctx.Guild.Roles.First(x => x.Value.Name == "Muted").Value;
            await member.GrantRoleAsync(role);
            await ctx.RespondAsync(":speaking_head:");
        }
    }
}
