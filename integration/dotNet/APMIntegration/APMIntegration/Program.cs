using System;
using System.Collections.Generic;
using System.IO.Pipes;
using System.Linq;
using System.Text;
using System.Threading;

namespace APMIntegration
{
    class Program
    {
        static void Main(string[] args)
        {
            TestClassInstance();
            //TestStaticMethod();
        }

        static void TestStaticMethod()
        {
            Random rnd = new Random();
            while (true)
            {
                GazerNamedPipe.Write("q_debug_pipe", "paramStatic1", (rnd.Next() % 100).ToString());
                GazerNamedPipe.Write("q_debug_pipe", "paramStatic2", (rnd.Next() % 100).ToString());
                GazerNamedPipe.Write("q_debug_pipe", "paramStatic3", (rnd.Next() % 100).ToString());
                Console.WriteLine("{0}", DateTime.Now.ToString());
                Thread.Sleep(100);
            }
        }

        static void TestClassInstance()
        {
            using (GazerNamedPipe gazerPipe = new GazerNamedPipe("q_debug_pipe"))
            {
                Random rnd = new Random();
                while (true)
                {
                    gazerPipe.Write("param1", (rnd.Next() % 1000).ToString());
                    gazerPipe.Write("param2", (rnd.Next() % 1000).ToString());
                    gazerPipe.Write("param3", (rnd.Next() % 1000).ToString());
                    Console.WriteLine("{0}", DateTime.Now.ToString());
                    Thread.Sleep(100);
                }
            }
        }
    }
}
