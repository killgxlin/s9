using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class PlayerStub : MonoBehaviour {
	public void Destroy() {
		Destroy (gameObject);
	}
	public void setPos(Msg.Vector2 pos) {
		transform.position = new Vector3 (pos.X, 0, pos.Y);
	}
	public void setCamPos(Msg.Vector2 pos) {
		transform.position = new Vector3 (pos.X, 10, pos.Y);
	}
}
