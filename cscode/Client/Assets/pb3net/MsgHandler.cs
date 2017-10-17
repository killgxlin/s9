using System;
using System.Collections.Generic;
using UnityEngine;
using Google.Protobuf;

/*

	class HandleObject {
		void On(object ctx, Msg.SCPong m) {
			Debug.Log (m);
		}
	}
*/

public class MsgHandler
{
	// single thread call -----------------------------------------------------
	object[] paras = new object[2];
	Type[] typs = new Type[2];
	object h;
	public void SetHandler(object h1) {
		h = h1;
	}
	public void HandleMessage(object ctx, object m) {
		var m1 = (IMessage)m;
		var ht = h.GetType ();

		typs [0] = ctx.GetType ();
		typs [1] = m1.GetType ();
		var mt = ht.GetMethod("On", System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Instance, null, typs, null);
		if (mt == null) {
			throw new Exception ("undefined handle for "+typs[0].ToString());
		}

		paras [0] = ctx;
		paras [1] = m1;
		mt.Invoke (h, paras);
	}

}
