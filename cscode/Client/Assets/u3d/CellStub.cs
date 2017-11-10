using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class CellStub : MonoBehaviour {
	public void Init(Msg.SEnterCell m) {
		var cell = m.Cell;
		Util.AddBorder (gameObject, cell.SwitchBorder, Color.red);
		Util.AddBorder (gameObject, cell.Border, Color.blue);
		Util.AddBorder (gameObject, cell.MirrorBorder, Color.green);

		for (var i = 0; i < m.Neighbor.Count; i++) {
			var n = m.Neighbor [i];
			Util.AddBorder (gameObject, n.SwitchBorder, Color.red);
			Util.AddBorder (gameObject, n.Border, Color.blue);
			Util.AddBorder (gameObject, n.MirrorBorder, Color.green);
		}
	}
	public void Destroy() {
		Destroy (gameObject);
	}
}
