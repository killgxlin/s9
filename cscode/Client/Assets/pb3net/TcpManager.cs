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

		public string currentHost{ get; set; }
		TcpClient tcp { get; set; }
		List<IMessage> msgRecved { get; set; }
		volatile OnDisconnect onDisconnect;

		public TcpManager()
		{
			msgRecved = new List<IMessage> ();
			tcp = new TcpClient ();
			//tcp.log = true; 调试log
			tcp.onDisconnect 	= Event_DisConnect;	//断线事件
			tcp.onMessage 		= Event_Message;	//接收消息
			tcp.encode 			= EnCode;			//加密
			tcp.decode 			= DeCode;			//解密

			onDisconnect = null;
		}

		public void Connect(string host, int port)
		{
			currentHost = host;
			tcp.Connect (host, port);
		}

		public void Disconnect()
		{
			if (tcp == null)
				return;
			tcp.DisConnect ();
		}

		void Event_DisConnect()
		{
			onDisconnect = OnDisConnect;
		}

		void OnDisConnect()
		{

		}

		void Event_Message( object message)
		{
			lock (msgRecved) msgRecved.Add (message as IMessage);
		}
			
		byte[] EnCode(object message)
		{
			return MsgCoder.ToByteArray(message);
		}
			
		object DeCode(byte[] packet)
		{
			return MsgCoder.ParseFrom (packet);
		}

		public void SendMessage(IMessage msg)
		{    
			if (tcp != null) {
				tcp.SendMessage (msg);
			}
		}

		public Pb3Net.NetState netState
		{
			get { return tcp==null ? Pb3Net.NetState.DisConnected :tcp.currentState; }
		}
			
		public List<IMessage> RecvMessage(int num) {
			if (msgRecved == null)
				return null;
			if (msgRecved.Count == 0)
				return null;
			lock (msgRecved) {
				var n = Math.Min (num, msgRecved.Count);
				var ret = msgRecved.GetRange (0, n);
				msgRecved.RemoveRange (0, n);
				return ret;
			}
		}

		public void Update(int fram, float dt)
		{
			//断线事件
			if(onDisconnect != null){
				onDisconnect.Invoke ();
				onDisconnect = null;
			}
		}
	}
}