using System;
using Google.Protobuf;
namespace Pb3Net
{
	public class MsgCoder
	{
		// message marshal and unmarshal ------------------------------------------
		static Google.Protobuf.Reflection.TypeRegistry cellType = Google.Protobuf.Reflection.TypeRegistry.FromFiles (Cell.ProtosReflection.Descriptor);

		public static byte[] ToByteArray(object m) {
			IMessage m1 = (IMessage)m;
			var msg = new Gate.Msg{
				Name = m1.Descriptor.FullName,
				Raw = m1.ToByteString (),
			};

			return msg.ToByteArray ();
		}

		public static object ParseFrom(byte[] raw) {
			var m = (Gate.Msg)Gate.Msg.Parser.ParseFrom (raw);
			return cellType.Find (m.Name).Parser.ParseFrom (m.Raw);;
		}

		public static void Test() {
//			var ping = new Cell.Ping{ Start=1, Content="hello"};
//			var b = ToByteArray (ping);
//			var c = ParseFrom (b);
//			var ping2 = c as Cell.Ping;
//			Console.WriteLine (ping2);
		}
	}
}

