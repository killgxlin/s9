using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class CellStub : MonoBehaviour {
	void addBorder(Msg.AABB box, Color color) {
		var r = gameObject.AddComponent<LineRenderer>();
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
		r.startColor = color;
		r.endColor = color;
	}
	public void Init(Msg.Cell cell) {
		addBorder (cell.SwitchBorder, Color.blue);
		//addBorder (cell.SwitchBorder, Color.red);
	}
	public void Destroy() {
		Destroy (gameObject);
	}
}
