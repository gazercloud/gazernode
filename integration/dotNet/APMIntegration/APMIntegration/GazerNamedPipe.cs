using System;
using System.IO.Pipes;
using System.Text;

namespace APMIntegration
{
    public class GazerNamedPipe : IDisposable
    {
        public GazerNamedPipe(string pipeName)
        {
            pipeName_ = pipeName;
        }

        private string pipeName_ = "";
        private NamedPipeClientStream stream_ = null;

        public void Write(string name, string value)
        {
            try
            {
                if (stream_ == null)
                {
                    stream_ = new NamedPipeClientStream(pipeName_);
                    stream_.Connect(0);
                }

                if (stream_ != null)
                {
                    var bytes = Encoding.UTF8.GetBytes(string.Format("{0}={1}\r\n", name, value));
                    stream_.Write(bytes, 0, bytes.Length);
                }
            }
            catch (Exception ex)
            {
                if (stream_ != null)
                {
                    stream_.Close();
                    stream_.Dispose();
                    stream_ = null;
                    //GC.Collect();
                }
                Console.WriteLine("GazerPipe Error: {0}", ex);
            }
        }

        public static void Write(string pipeName, string name, string value)
        {
            using (NamedPipeClientStream stream = new NamedPipeClientStream(pipeName))
            {
                try
                {
                    stream.Connect(0);
                    var bytes = Encoding.UTF8.GetBytes(string.Format("{0}={1}\r\n", name, value));
                    stream.Write(bytes, 0, bytes.Length);
                    stream.Close();
                }
                catch (Exception ex)
                {
                    //Console.WriteLine("GazerPipe Error: {0}", ex);
                }
            }

            //GC.Collect();
        }

        public void Dispose()
        {
            if (stream_ != null)
            {
                stream_.Close();
                stream_.Dispose();
                stream_ = null;
                //GC.Collect();
            }
        }
    }
}
