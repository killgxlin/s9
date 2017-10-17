using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;


namespace Pb3Net
{
	public enum NetState
	{
		DisConnected,
		Connecting,
		Connected
	}

	public delegate void OnMessage(object message);
	public delegate void OnDisconnect();
	public delegate byte[] Encode(object message);
	public delegate object Decode(byte[] message);

	public class TcpClient : IDisposable
	{
		string host { get; set; }
		int port { get; set; }
		IPEndPoint address { get; set; }
		Socket socket { get; set; }
		volatile NetState state;
		public NetState currentState { get { return state; } }
		public bool log { get; set; }

		DataStream recvStream { get; set; }
		DataStream buffStream { get; set; }

		public OnMessage onMessage { get; set; }
		public OnDisconnect onDisconnect { get; set; }
		public Encode encode { get; set; }
		public Decode decode { get; set; }

		public TcpClient(string host, int port)
		{
			this.host = host;
			this.port = port;

			Init();
		}

		public TcpClient()
		{
			Init();
		}

		void Init()
		{
			this.state = NetState.DisConnected;
			this.recvStream = new DataStream(1024 * 64);
			this.buffStream = new DataStream(1024 * 64);
		}

		public void Dispose()
		{
			state = NetState.DisConnected;
			host = string.Empty;
			port = 0;
			buffStream.Clear();
			recvStream.Clear();
		}

		public void Connect()
		{
			Connect(host, port);
		}

		public void Connect(string host, int port)
		{
			try
			{
				this.host = host;
				this.port = port;
				this.address = new IPEndPoint(IPAddress.Parse(host), port);

				//初始化socket
				socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
				//更新状态
				state = NetState.Connecting;
				//开始链接主机
				socket.BeginConnect(address, new AsyncCallback(EndConnect), null);
			}
			catch (Exception ex)
			{
				PrintLog("connect to {1} error! {0}", ex.Message, address.ToString());

				//断开事件 
				DisconnectEvent();
			}
		}

		void EndConnect(IAsyncResult result)
		{
			try
			{
				//链接完成
				socket.EndConnect(result);
				//更新状态
				state = NetState.Connected;
				//开始接收消息
				RecvBeginProcess();
			}
			catch (Exception ex)
			{ 
				PrintLog("connect to {1} end error! {0}", ex.Message, address.ToString());

				//断开事件 
				DisconnectEvent();
			}
		}

		public void DisConnect()
		{
			try
			{
				if (state == NetState.DisConnected) return;
				if (!socket.Connected) return;

				//关闭收发
				socket.Shutdown(SocketShutdown.Both);
				//开始尝试断开链接
				socket.BeginDisconnect(false, new AsyncCallback(DisconnectEnd), null);
			}
			catch (Exception ex)
			{
				PrintLog("disconnect to {1} error! {0}", ex.Message, address.ToString());

				//断开事件 
				DisconnectEvent();

				//清理数据
				Dispose();
			}
		}

		void DisconnectEnd(IAsyncResult result)
		{
			try
			{
				//结束断开链接
				socket.EndDisconnect(result);

				//断开事件 
				DisconnectEvent();

				//清理数据
				Dispose();
			}
			catch (Exception ex)
			{ 
				PrintLog("disconnect to {1} end error! {0}", ex.Message, address.ToString());

				//断开事件 
				DisconnectEvent();

				//清理数据
				Dispose();
			}
		}

		void DisconnectEvent()
		{
			if (state == NetState.DisConnected)
				return;
			state = NetState.DisConnected;

			if (onDisconnect != null) onDisconnect.Invoke();
		}

		void RecvBeginProcess()
		{
			try
			{
				if (socket == null) return;
				if (state != NetState.Connected) return;

				socket.BeginReceive(
					recvStream.buff,
					0,
					recvStream.buff.Length,
					SocketFlags.None,
					new AsyncCallback(RecvEndProcess),
					null
				);
			}
			catch (Exception ex)
			{ 
				PrintLog("recv message error! {0}", ex.Message);

				//断开事件 
				DisconnectEvent();

				//清理数据
				Dispose ();
			}
		}

		void RecvEndProcess(IAsyncResult result)
		{
			try
			{
				//结束监听数据
				int length = socket.EndReceive(result);

				if (length <= 0) {
					//断开事件 
					DisconnectEvent();

					return;
				}
				//copy到缓冲区
				buffStream.WriteBytes(recvStream.buff, length);
				recvStream.Clear();
				//处理数据
				ProcessPacket();

				//继续监听数据
				RecvBeginProcess();
			}
			catch (Exception ex)
			{
				PrintLog("recv message end error! {0}", ex.Message);

				//断开事件 
				DisconnectEvent();

				//清理数据
				Dispose ();
			}
		}

		void ProcessPacket()
		{
			int offset = buffStream.position;   						//起始位置
			int length = buffStream.length + buffStream.position;     	//数据长度

			//解析所有封包
			while (offset < length){

				//获取网络封包
				var raw = CreateNetPacket(buffStream.buff, offset, length);

				//封包无效的跳出
				if (raw == null) break;

				//缓冲区位置位移
				offset += (raw.Length + 4);

				//解密封包
				if (onMessage != null) onMessage.Invoke(decode(raw));
			}

			//清理封包
			if (offset >= length) {
				buffStream.Clear();
				return;
			}

			if (offset == 0) return;
				
			//封包里还有残留数据，挪动数据
			Buffer.BlockCopy(buffStream.buff, offset, buffStream.buff, 0, length-offset);
			buffStream.position = 0;
			buffStream.length = length - offset;
		}	

		public void SendMessage(object message)
		{
			if (state != NetState.Connected) return;

			//开始发送数据
			SendBeginProcess(CreateSendBytes(encode(message)));
		}

		void SendBeginProcess(byte[] bytes)
		{
			try
			{
				if (socket == null) return;

				socket.BeginSend(bytes, 0, bytes.Length, SocketFlags.None, new AsyncCallback(SendEndProcess), this);
			}
			catch (Exception ex)
			{
				PrintLog("send message error! {0}", ex.Message);
			}
		}

		void SendEndProcess(IAsyncResult result)
		{
			try
			{
				if (socket == null) return;
				socket.EndSend(result);
			}
			catch (Exception ex)
			{
				PrintLog("send message end error! {0}", ex.Message);
			}

		}

		////////////////////////////////////////////////////////////////////////////////////////////////////////
		/// 工具函数
		void PrintLog(string message, params object[] args)
		{
			if (!log) return;

			Console.Write("[{0}]{1}", Thread.CurrentThread.ManagedThreadId, string.Format(message, args));
		}



		////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
		/// 组合封包,分解封包

		byte[] CreateSendBytes(byte[] body)
		{
			byte[] buff = new byte[4 + body.Length];

			Buffer.BlockCopy(BitConverter.GetBytes(IPAddress.HostToNetworkOrder(buff.Length-4)), 0, buff, 0, 4);
			Buffer.BlockCopy(body, 0, buff, 4, body.Length);

			return buff;
		}

		byte[] CreateNetPacket(byte[] buff, int offset, int length)
		{
			int current = offset + 4;

			//不满足包头长度的话
			if (current > length) return null;

			//获取总长度
			int buffLen = IPAddress.NetworkToHostOrder(BitConverter.ToInt32(buff, offset));

			//包体总长度超过缓冲区长度
			if ((buffLen + current) > length) return null;

			var raw = new byte[buffLen];

			Buffer.BlockCopy(buff, offset + 4, raw, 0, raw.Length);

			return raw;
		}
	}
}

