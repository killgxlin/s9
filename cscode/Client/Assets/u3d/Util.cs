using System;
using UnityEngine;

public class Util
{
	static public bool Approximately (float a, float b)
	{
		return UnityEngine.Mathf.Approximately (a, b);
	}
	static public PlayerStub AttachStub(string name) {
		var obj = GameObject.Find (name);
		if (obj == null)
			return null;

		var stub = obj.GetComponent<PlayerStub> ();
		if (stub == null)
			stub = obj.AddComponent<PlayerStub> ();

		return stub;
	}
	static public PlayerStub CreateStub(Msg.PlayerData data) {
		var obj = GameObject.CreatePrimitive (PrimitiveType.Cylinder);
		obj.name = string.Format ("gameobj_{0}", data.Id);

		var stub = obj.AddComponent<PlayerStub> ();
		stub.setPos (data.Pos);

		return stub;
	}
	static public CellStub CreateCell(Msg.Cell cell) {
		var obj = new GameObject (cell.Name);
		var stub = obj.AddComponent<CellStub>();
		stub.Init (cell);
		return stub;
	}

	static public Msg.Vector2 GetInput() {
		var dir = Vector3.zero;
		if (Input.GetKey (KeyCode.A)) {
			dir.x = -1;
		}
		if (Input.GetKey (KeyCode.D)) {
			dir.x = 1;
		}
		if (Input.GetKey (KeyCode.W)) {
			dir.z = 1;
		}
		if (Input.GetKey (KeyCode.S)) {
			dir.z = -1;
		}
		if (dir != Vector3.zero)
			dir.Normalize ();
		
		return new Msg.Vector2{ X = dir.x, Y = dir.z };
	}
}
