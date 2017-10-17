using System;
namespace Pb3Net
{
	internal class DataStream
	{
		internal byte[] buff;
		internal int length;
		internal int position;

		internal DataStream(int len)
		{
			position = length = 0;
			buff = new byte[len];
		}

		internal void Clear()
		{ 
			position = length = 0;
		}

		internal void WriteBytes(byte[] array, int len)
		{
			int real = len + (position + length);

			if(real > buff.Length)
				Array.Resize<byte>(ref buff, real);

			Buffer.BlockCopy(array, 0, buff, length, len);

			length += len;
		}
	}
}

