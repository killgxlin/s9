using System.Collections;
using System.Collections.Generic;

public class Game {
	static public Game Instance;
	MsgHandler handler = new MsgHandler ();
	Pb3Net.TcpManager conn = new Pb3Net.TcpManager();
	PlayerData self;
	Stub cam;
	public void Start (Main main) {
		Instance = this;

		handler.SetHandler (this);
		//conn.Connect ("10.235.226.80", 9000);
		conn.Connect ("192.168.31.233", 9000);
		cam = Util.AttachStub ("Main Camera");
	}
		
	public void Update (Main main, float deltaTime) {
		updateMessage ();
		updatePlayers (deltaTime);

		if (self != null) {
			updateCamera ();
			updateInput ();
		}
	}
	/// <summary>
	/// player
	/// </summary>
	public class PlayerData {
		public Stub stub;
		public Msg.PlayerData data;
	}
	Dictionary<int, PlayerData> players = new Dictionary<int, PlayerData>();
	void updatePlayers(float delta) {
		for (var itr = players.GetEnumerator (); itr.MoveNext ();) {
			var data = itr.Current.Value.data;
			if (Util.Approximately (data.Vel.X, 0) && Util.Approximately (data.Vel.Y, 0))
				continue;
			
			data.Pos.X += data.Vel.X * delta;
			data.Pos.Y += data.Vel.Y * delta;
			itr.Current.Value.stub.setPos (data.Pos);
		}
	}

	/// <summary>
	/// camera
	/// </summary>
	void updateCamera ()
	{
		var pos = self.data.Pos.Clone ();
		pos.Y = 10;
		cam.setPos (pos);
	}

	/// <summary>
	/// input
	/// </summary>
	Msg.Vector2 lastDir = new Msg.Vector2();
	void updateInput ()
	{
		var dir = Util.GetInput ();
		if (!lastDir.Equals (dir)) {
			var d = self.data.Clone ();
			d.Vel = dir;
			conn.SendMessage (new Msg.CUpdate {
				Data = d
			});
		}
		lastDir = dir.Clone ();
	}

	/// <summary>
	/// message
	/// </summary>
	Pb3Net.NetState lastSt = Pb3Net.NetState.Connecting;
	void updateMessage() {
		if (lastSt == Pb3Net.NetState.Connecting && conn.netState == Pb3Net.NetState.Connected) {
			conn.SendMessage (new Msg.CLogin{Account="hello"});
		}
		lastSt = conn.netState;
		var msgs = conn.RecvMessage (100);
		if (msgs != null) {
			for (var itr = msgs.GetEnumerator (); itr.MoveNext ();) {
				handler.HandleMessage (this, itr.Current);
			}
		}
	}
	void On(object ctx, Msg.SEnterCell m) {
		var stub = Util.CreateStub (m.Self);
		self = new PlayerData{ data = m.Self, stub = stub };
		players.Add (m.Self.Id, self);

	}
	void On(object ctx, Msg.SLeaveCell m) {
		//m.CalculateSize;
		var p = players[self.data.Id];
		if (p == null)
			return;

		p.stub.Destroy ();
		players.Remove (p.data.Id);
	}
	void On(object ctx, Msg.SAdd m) {
		for (var itr = m.Data.GetEnumerator (); itr.MoveNext ();) {
			var stub = Util.CreateStub (itr.Current);
			players.Add (itr.Current.Id, new PlayerData{ data = itr.Current, stub = stub });
		}
	}
	void On(object ctx, Msg.SRemove m) {

		for (var itr = m.Id.GetEnumerator (); itr.MoveNext ();) {
			PlayerData data;
			if (!players.TryGetValue (itr.Current, out data))
				continue;
			data.stub.Destroy ();
			players.Remove (itr.Current);
		}
	}
	void On(object ctx, Msg.SUpdate m) {
		PlayerData data;
		if (!players.TryGetValue (m.Data.Id, out data))
			return;
		data.data = m.Data;
	}
}
