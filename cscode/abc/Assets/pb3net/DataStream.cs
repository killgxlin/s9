using System;
namespace Pb3Net
{
	/// <summary>
	/// 数据流
	/// </summary>
	internal class DataStream
	{
		/// <summary>
		/// 数据缓冲区
		/// </summary>
		/// <value>The buff.</value>
		internal byte[] buff;

		/// <summary>
		/// 数据缓冲区当前长度
		/// </summary>
		/// <value>The current Lenght.</value>
		internal int length;

		/// <summary>
		/// 当前缓冲区起始位置
		/// </summary>
		/// <value>The position.</value>
		internal int position;

		/// <summary>
		/// 初始化大小
		/// </summary>
		internal DataStream(int len)
		{
			position = length = 0;
			buff = new byte[len];
		}

		/// <summary>
		/// 清理数据
		/// </summary>
		internal void Clear()
		{ 
			position = length = 0;
		}

		/// <summary>
		/// 写入缓冲区
		/// </summary>
		/// <returns>The bytes.</returns>
		/// <param name="array">Array.</param>
		/// <param name="len">Length.</param>
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

