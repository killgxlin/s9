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
	static public CellStub CreateCell(Msg.SEnterCell m) {
		var obj = new GameObject (m.Cell.Name);
		var stub = obj.AddComponent<CellStub>();
		stub.Init (m);
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

		//dir.Scale (new Vector3 (0.5f, 0.5f, 0.5f));
		
		return new Msg.Vector2{ X = dir.x, Y = dir.z };
	}

	static public void AddBorder(GameObject gameObject, Msg.AABB box, Color color) {
		var g2 = new GameObject ();
		var r = g2.AddComponent<LineRenderer> ();
		g2.transform.SetParent (gameObject.transform);
		var poses = new Vector3[] {
			new Vector3 (box.Minx, 0, box.Miny),
			new Vector3 (box.Minx, 0, box.Maxy),
			new Vector3 (box.Maxx, 0, box.Maxy),
			new Vector3 (box.Maxx, 0, box.Miny),
			new Vector3 (box.Minx, 0, box.Miny),
		};

		r.numPositions = poses.Length;
		r.SetPositions (poses);
		r.startWidth = 0.1f;
		r.endWidth = 0.1f;
		var mat = new Material (Shader.Find ("Legacy Shaders/Diffuse"));
		mat.color = color;
		r.material = mat;
	}
}
