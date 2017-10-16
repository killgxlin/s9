using UnityEngine;
using System;
using System.Collections.Generic;
using Google.Protobuf;


namespace Pb3Net
{


	/// <summary>
	/// tcp链接管理
	/// </summary>
	public class TcpManager //: Singleton<TcpManager>
	{
		/// <summary>
		/// 当前链接的服务器
		/// </summary>
		/// <value>The current host.</value>
		public string currentHost{ get; set; }

		/// <summary>
		/// tcp链接
		/// </summary>
		/// <value>The tcp.</value>
		TcpClient tcp { get; set; }


		/// <summary>
		/// 消息队列
		/// </summary>
		/// <value>The tcp queue.</value>
		Queue<IMessage> tcpQueue { get; set; }

		/// <summary>
		/// 断线事件
		/// </summary>
		volatile OnDisconnect onDisconnect;

		/// <summary>
		/// 构造函数，创建内部对象
		/// </summary>
		public TcpManager()
		{
			tcpQueue = new Queue<IMessage> ();
			tcp = new TcpClient ();
			//tcp.log = true; 调试log
			tcp.onDisconnect 	= Event_DisConnect;	//断线事件
			tcp.onMessage 		= Event_Message;	//接收消息
			tcp.encode 			= EnCode;			//加密
			tcp.decode 			= DeCode;			//解密

			onDisconnect = null;
		}

		/// <summary>
		/// 连接到服务器
		/// </summary>
		/// <param name="hostIP">主机IP地址</param>
		/// <param name="port">端口</param>
		public void Connect(string host, int port)
		{
			currentHost = host;
			tcp.Connect (host, port);
		}

		/// <summary>
		/// 断开链接
		/// </summary>
		public void Disconnect()
		{
			if (tcp == null)
				return;
			tcp.DisConnect ();
		}

		/// <summary>
		/// 断开链接事件
		/// 此事件为异步调用
		/// </summary>
		void Event_DisConnect()
		{
			onDisconnect = OnDisConnect;
		}

		/// <summary>
		/// 断开事件 
		/// 同步调用
		/// </summary>
		void OnDisConnect()
		{

		}

		/// <summary>
		/// 接收到事件
		/// 此事件为异步调用
		/// </summary>
		/// <param name="msgId">Message identifier.</param>
		/// <param name="message">Message.</param>
		void Event_Message( object message)
		{
			lock (tcpQueue) tcpQueue.Enqueue (message as IMessage);
		}

		/// <summary>
		/// 加密函数
		/// </summary>
		/// <returns>The code.</returns>
		/// <param name="message">Message.</param>
		byte[] EnCode(object message)
		{
			return MsgCoder.ToByteArray(message);
		}

		/// <summary>
		/// 解密封包
		/// </summary>
		/// <returns>The code.</returns>
		/// <param name="packet">Packet.</param>
		object DeCode(byte[] packet)
		{
			return MsgCoder.ParseFrom (packet);
		}

		/// <summary>
		/// 发消息
		/// </summary>
		/// <param name="msgID">消息ID</param>
		/// <param name="msgBuilder">外部构造好的消息</param>
		public void SendMessage(IMessage msg)
		{    
			if (tcp != null) {
				tcp.SendMessage (msg);
			
			}
		}

		/// <summary>
		/// 当前网络链接状态
		/// </summary>
		public Pb3Net.NetState netState
		{
				get { return tcp==null ? Pb3Net.NetState.DisConnected :tcp.currentState; }
		}


		/// <summary>
		/// 执行封包
		/// </summary>
		/// <param name="packet">Packet.</param>
		void RunMsgHandler(IMessage message)
		{
			MsgHandler.HandleMessage (this, message);
		}

		/// <summary>
		/// 封包处理
		/// </summary>
		void PacketProcess()
		{
			if (tcpQueue == null)
				return;
			if (tcpQueue.Count == 0)
				return;
			IMessage[] array = null;

			lock (tcpQueue) {
				array = tcpQueue.ToArray ();
				tcpQueue.Clear ();
			}

			for (int i = 0; i < array.Length; i++) {
				RunMsgHandler (array[i]);
			}
		}

		/// <summary>
		/// 在游戏循环中调用网络处理
		/// </summary>
		public void Update(int fram, float dt)
		{
			//封包处理
			PacketProcess ();

			//断线事件
			if(onDisconnect != null){
				onDisconnect.Invoke ();
				onDisconnect = null;
			}
		}
	}
}