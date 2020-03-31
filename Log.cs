using System;

namespace ThomasBot
{
    class Log
    {
        /// <summary>
        /// Writes a log message for the given log level.
        /// </summary>
        /// <param name="message">The text to write.</param>
        /// <param name="level">The  error-level.</param>
        public static void WriteLogMessage(string message, LogOutputLevel level)
        {
#if DEBUG
            if (level == LogOutputLevel.Debug)
            {
                Console.WriteLine($"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] [Log] [Debug] {message}");
                return;
            }
#endif

            switch (level)
            {
                case LogOutputLevel.Info:
                    WriteColor($"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] [Log] [Info] {message}", ConsoleColor.White);
                    break;
                case LogOutputLevel.Warning:
                    WriteColor($"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] [Log] [Warning] {message}", ConsoleColor.Yellow);
                    break;
                case LogOutputLevel.Error:
                    WriteColor($"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] [Log] [Error] {message}", ConsoleColor.Red);
                    break;
                case LogOutputLevel.Critical:
                    WriteColor($"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] [Log] [Critical] {message}", ConsoleColor.Red);
                    break;
                case LogOutputLevel.Good:
                    WriteColor($"[{DateTime.Now.ToLocalTime().ToString("yyyy-MM-dd HH:mm:ss zzz")}] [Log] [Good] {message}", ConsoleColor.Green);
                    break;
                default:
                    break;
            }

        }

        private static void WriteColor(string message, ConsoleColor color)
        {
            var oldColor = Console.ForegroundColor;
            Console.ForegroundColor = color;
            Console.WriteLine(message);
            Console.ForegroundColor = oldColor;
        }
    }


    enum LogOutputLevel
    {
        Debug = 0,
        Info = 1,
        Warning = 2,
        Error = 3,
        Critical = 4,
        Good = 5
    }
}
