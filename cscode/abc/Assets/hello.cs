using System.Collections;
using System.Collections.Generic;
using UnityEngine;


public class hello : MonoBehaviour {
	
	Pb3Net.TcpManager man = new Pb3Net.TcpManager();
	Cell.PlayerData self;
	// Use this for initialization
	void Start () {
		MsgHandler.SetHandler (this);
		man.Connect ("10.235.226.80", 9000);
		lastSt = man.netState;
	}

	Pb3Net.NetState lastSt;
	// Update is called once per frame
	void Update () {
		if (man.netState == Pb3Net.NetState.Connected && self!=null) {
			self.Pos.X = 1.0f;
			self.Pos.Y = 2.0f;
			self.Pos.Z = 3.0f;
			self.Vel.X = 11.0f;
			self.Vel.Y = 21.0f;
			self.Vel.Z = 31.0f;
			man.SendMessage (new Cell.CMove{Data=self});
			self = null;
		}
		lastSt = man.netState;

		man.Update (1, Time.deltaTime);
	}

	// message handlers -----------------------------------------------------
	void On(object ctx, Cell.SEnter m) {
		Debug.Log (m);
		self = m.Self;
	}
	void On(object ctx, Cell.SAdd m) {
		Debug.Log (m);
	}
	void On(object ctx, Cell.SRemove m) {
		Debug.Log (m);
	}
	void On(object ctx, Cell.SMove m) {
		Debug.Log (m);
	}
}
