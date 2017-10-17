using System;
using UnityEngine;

public class Util
{
	static public bool Approximately (float a, float b)
	{
		return UnityEngine.Mathf.Approximately (a, b);
	}
	static public Stub AttachStub(string name) {
		var obj = GameObject.Find (name);
		if (obj == null)
			return null;

		var stub = obj.GetComponent<Stub> ();
		if (stub == null)
			stub = obj.AddComponent<Stub> ();

		return stub;
	}
	static public Stub CreateStub(Cell.PlayerData data) {
		var obj = GameObject.CreatePrimitive (PrimitiveType.Cylinder);
		obj.name = string.Format ("gameobj_{0}", data.Id);

		var stub = obj.AddComponent<Stub> ();
		stub.setPos (data.Pos);

		return stub;
	}
	static public Cell.Vector3 GetInput() {
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
		
		return new Cell.Vector3{ X = dir.x, Y = dir.y, Z = dir.z };
	}
}
