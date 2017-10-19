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
		conn.Connect ("10.235.226.80", 9000);
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
			if (Util.Approximately (data.Vel.X, 0) && Util.Approximately (data.Vel.Y, 0) && Util.Approximately (data.Vel.Z, 0))
				continue;
			
			data.Pos.X += data.Vel.X * delta;
			data.Pos.Y += data.Vel.Y * delta;
			data.Pos.Z += data.Vel.Z * delta;
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
	Msg.Vector3 lastDir = new Msg.Vector3();
	void updateInput ()
	{
		var dir = Util.GetInput ();
		if (!lastDir.Equals (dir)) {
			var d = self.data.Clone ();
			d.Vel = dir;
			conn.SendMessage (new Msg.CMove {
				Data = d
			});
		}
		lastDir = dir.Clone ();
	}

	/// <summary>
	/// message
	/// </summary>
	void updateMessage() {
		var msgs = conn.RecvMessage (100);
		if (msgs != null) {
			for (var itr = msgs.GetEnumerator (); itr.MoveNext ();) {
				handler.HandleMessage (this, itr.Current);
			}
		}
	}
	void On(object ctx, Msg.SEnter m) {
		var stub = Util.CreateStub (m.Self);
		self = new PlayerData{ data = m.Self, stub = stub };
		players.Add (m.Self.Id, self);
		for (var itr = m.Other.GetEnumerator (); itr.MoveNext ();) {
			stub = Util.CreateStub (itr.Current);
			players.Add (itr.Current.Id, new PlayerData{ data = itr.Current, stub = stub });
		}
	}
	void On(object ctx, Msg.SAdd m) {
		var sub = Util.CreateStub (m.Data);
		players.Add (m.Data.Id, new PlayerData{ data = m.Data, stub = sub });
	}
	void On(object ctx, Msg.SRemove m) {
		PlayerData data;
		if (!players.TryGetValue (m.Id, out data))
			return;
		data.stub.Destroy ();
		players.Remove (m.Id);
	}
	void On(object ctx, Msg.SMove m) {
		PlayerData data;
		if (!players.TryGetValue (m.Data.Id, out data))
			return;
		data.data = m.Data;
	}
}
