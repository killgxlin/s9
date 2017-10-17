using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Main : MonoBehaviour {
	Game game = new Game();
	// Use this for initialization
	void Start () {
		game.Start (this);
	}
	
	// Update is called once per frame
	void Update () {
		game.Update (this, Time.deltaTime);
	}


}
