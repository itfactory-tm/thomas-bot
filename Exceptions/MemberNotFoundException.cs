using System;

namespace ThomasBot.Exceptions
{
    class MemberNotFoundException : Exception
    {
        public MemberNotFoundException(string message) : base(message)
        {

        }
    }
}
