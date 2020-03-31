using System.Threading.Tasks;

namespace ThomasBot
{
    class Program
    {
        static async Task Main()
        {
            var tb = new ThomasBot();
            await tb.RunBotAsync();
        }
    }
}
