using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Stub : MonoBehaviour {
	public void Destroy() {
		Destroy (gameObject);
	}
	public void setPos(Cell.Vector3 pos) {
		transform.position = new Vector3 (pos.X, pos.Y, pos.Z);
	}
}
