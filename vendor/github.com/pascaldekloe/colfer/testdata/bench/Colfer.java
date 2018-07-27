package testdata.bench;


// This file was generated by colf(1); DO NOT EDIT


import static java.lang.String.format;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.util.InputMismatchException;
import java.nio.BufferOverflowException;
import java.nio.BufferUnderflowException;


/**
 * Data bean with built-in serialization support.
 * @author generated by colf(1)
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class Colfer implements java.io.Serializable {

	/** The upper limit for serial byte sizes. */
	public static int colferSizeMax = 16 * 1024 * 1024;


	public long key;
	public String host = "";
	public int port;
	public long size;
	public long hash;
	public double ratio;
	public boolean route;


	/**
	 * {@link #reset(InputStream) Reusable} deserialization of Colfer streams.
	 */
	public static class Unmarshaller {

		/** The data source. */
		protected InputStream in;

		/** The read buffer. */
		protected byte[] buf;

		/** The {@link #buf buffer}'s data start index, inclusive. */
		protected int offset;

		/** The {@link #buf buffer}'s data end index, exclusive. */
		protected int i;


		/**
		 * @param in the data source or {@code null}.
		 * @param buf the initial buffer or {@code null}.
		 */
		public Unmarshaller(InputStream in, byte[] buf) {
			if (buf == null || buf.length == 0)
				buf = new byte[Math.min(Colfer.colferSizeMax, 2048)];
			this.buf = buf;
			reset(in);
		}

		/**
		 * Reuses the marshaller.
		 * @param in the data source or {@code null}.
		 * @throws IllegalStateException on pending data.
		 */
		public void reset(InputStream in) {
			if (this.i != this.offset) throw new IllegalStateException("colfer: pending data");
			this.in = in;
			this.offset = 0;
			this.i = 0;
		}

		/**
		 * Deserializes the following object.
		 * @return the result or {@code null} when EOF.
		 * @throws IOException from the input stream.
		 * @throws SecurityException on an upper limit breach defined by {@link #colferSizeMax}.
		 * @throws InputMismatchException when the data does not match this object's schema.
		 */
		public Colfer next() throws IOException {
			if (in == null) return null;

			while (true) {
				if (this.i > this.offset) {
					try {
						Colfer o = new Colfer();
						this.offset = o.unmarshal(this.buf, this.offset, this.i);
						return o;
					} catch (BufferUnderflowException e) {
					}
				}
				// not enough data

				if (this.i <= this.offset) {
					this.offset = 0;
					this.i = 0;
				} else if (i == buf.length) {
					byte[] src = this.buf;
					if (offset == 0) this.buf = new byte[Math.min(Colfer.colferSizeMax, this.buf.length * 4)];
					System.arraycopy(src, this.offset, this.buf, 0, this.i - this.offset);
					this.i -= this.offset;
					this.offset = 0;
				}
				assert this.i < this.buf.length;

				int n = in.read(buf, i, buf.length - i);
				if (n < 0) {
					if (this.i > this.offset)
						throw new InputMismatchException("colfer: pending data with EOF");
					return null;
				}
				assert n > 0;
				i += n;
			}
		}

	}


	/**
	 * Serializes the object.
	 * @param out the data destination.
	 * @param buf the initial buffer or {@code null}.
	 * @return the final buffer. When the serial fits into {@code buf} then the return is {@code buf}.
	 *  Otherwise the return is a new buffer, large enough to hold the whole serial.
	 * @throws IOException from {@code out}.
	 * @throws IllegalStateException on an upper limit breach defined by {@link #colferSizeMax}.
	 */
	public byte[] marshal(OutputStream out, byte[] buf) throws IOException {
		if (buf == null || buf.length == 0)
			buf = new byte[Math.min(Colfer.colferSizeMax, 2048)];

		while (true) {
			int i;
			try {
				i = marshal(buf, 0);
			} catch (BufferOverflowException e) {
				buf = new byte[Math.min(Colfer.colferSizeMax, buf.length * 4)];
				continue;
			}

			out.write(buf, 0, i);
			return buf;
		}
	}

	/**
	 * Serializes the object.
	 * @param buf the data destination.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferOverflowException when {@code buf} is too small.
	 * @throws IllegalStateException on an upper limit breach defined by {@link #colferSizeMax}.
	 */
	public int marshal(byte[] buf, int offset) {
		int i = offset;

		try {
			if (this.key != 0) {
				long x = this.key;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (0 | 0x80);
				} else
					buf[i++] = (byte) 0;
				for (int n = 0; n < 8 && (x & ~0x7fL) != 0; n++) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}

			if (! this.host.isEmpty()) {
				buf[i++] = (byte) 1;
				int start = ++i;

				String s = this.host;
				for (int sIndex = 0, sLength = s.length(); sIndex < sLength; sIndex++) {
					char c = s.charAt(sIndex);
					if (c < '\u0080') {
						buf[i++] = (byte) c;
					} else if (c < '\u0800') {
						buf[i++] = (byte) (192 | c >>> 6);
						buf[i++] = (byte) (128 | c & 63);
					} else if (c < '\ud800' || c > '\udfff') {
						buf[i++] = (byte) (224 | c >>> 12);
						buf[i++] = (byte) (128 | c >>> 6 & 63);
						buf[i++] = (byte) (128 | c & 63);
					} else {
						int cp = 0;
						if (++sIndex < sLength) cp = Character.toCodePoint(c, s.charAt(sIndex));
						if ((cp >= 1 << 16) && (cp < 1 << 21)) {
							buf[i++] = (byte) (240 | cp >>> 18);
							buf[i++] = (byte) (128 | cp >>> 12 & 63);
							buf[i++] = (byte) (128 | cp >>> 6 & 63);
							buf[i++] = (byte) (128 | cp & 63);
						} else
							buf[i++] = (byte) '?';
					}
				}
				int size = i - start;
				if (size > Colfer.colferSizeMax)
					throw new IllegalStateException(format("colfer: field testdata/bench.Colfer.host size %d exceeds %d UTF-8 bytes", size, Colfer.colferSizeMax));

				int ii = start - 1;
				if (size > 0x7f) {
					i++;
					for (int x = size; x >= 1 << 14; x >>>= 7) i++;
					System.arraycopy(buf, start, buf, i - size, size);

					do {
						buf[ii++] = (byte) (size | 0x80);
						size >>>= 7;
					} while (size > 0x7f);
				}
				buf[ii] = (byte) size;
			}

			if (this.port != 0) {
				int x = this.port;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (2 | 0x80);
				} else
					buf[i++] = (byte) 2;
				while ((x & ~0x7f) != 0) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}

			if (this.size != 0) {
				long x = this.size;
				if (x < 0) {
					x = -x;
					buf[i++] = (byte) (3 | 0x80);
				} else
					buf[i++] = (byte) 3;
				for (int n = 0; n < 8 && (x & ~0x7fL) != 0; n++) {
					buf[i++] = (byte) (x | 0x80);
					x >>>= 7;
				}
				buf[i++] = (byte) x;
			}

			if (this.hash != 0) {
				long x = this.hash;
				if ((x & ~((1 << 49) - 1)) != 0) {
					buf[i++] = (byte) (4 | 0x80);
					buf[i++] = (byte) (x >>> 56);
					buf[i++] = (byte) (x >>> 48);
					buf[i++] = (byte) (x >>> 40);
					buf[i++] = (byte) (x >>> 32);
					buf[i++] = (byte) (x >>> 24);
					buf[i++] = (byte) (x >>> 16);
					buf[i++] = (byte) (x >>> 8);
					buf[i++] = (byte) (x);
				} else {
					buf[i++] = (byte) 4;
					while (x > 0x7fL) {
						buf[i++] = (byte) (x | 0x80);
						x >>>= 7;
					}
					buf[i++] = (byte) x;
				}
			}

			if (this.ratio != 0.0) {
				buf[i++] = (byte) 5;
				long x = Double.doubleToRawLongBits(this.ratio);
				buf[i++] = (byte) (x >>> 56);
				buf[i++] = (byte) (x >>> 48);
				buf[i++] = (byte) (x >>> 40);
				buf[i++] = (byte) (x >>> 32);
				buf[i++] = (byte) (x >>> 24);
				buf[i++] = (byte) (x >>> 16);
				buf[i++] = (byte) (x >>> 8);
				buf[i++] = (byte) (x);
			}

			if (this.route) {
				buf[i++] = (byte) 6;
			}

			buf[i++] = (byte) 0x7f;
			return i;
		} catch (ArrayIndexOutOfBoundsException e) {
			if (i - offset > Colfer.colferSizeMax)
				throw new IllegalStateException(format("colfer: serial exceeds %d bytes", Colfer.colferSizeMax));
			if (i > buf.length) throw new BufferOverflowException();
			throw e;
		}
	}

	/**
	 * Deserializes the object.
	 * @param buf the data source.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferUnderflowException when {@code buf} is incomplete. (EOF)
	 * @throws SecurityException on an upper limit breach defined by {@link #colferSizeMax}.
	 * @throws InputMismatchException when the data does not match this object's schema.
	 */
	public int unmarshal(byte[] buf, int offset) {
		return unmarshal(buf, offset, buf.length);
	}

	/**
	 * Deserializes the object.
	 * @param buf the data source.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @param end the index limit for {@code buf}, exclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferUnderflowException when {@code buf} is incomplete. (EOF)
	 * @throws SecurityException on an upper limit breach defined by {@link #colferSizeMax}.
	 * @throws InputMismatchException when the data does not match this object's schema.
	 */
	public int unmarshal(byte[] buf, int offset, int end) {
		if (end > buf.length) end = buf.length;
		int i = offset;

		try {
			byte header = buf[i++];

			if (header == (byte) 0) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.key = x;
				header = buf[i++];
			} else if (header == (byte) (0 | 0x80)) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.key = -x;
				header = buf[i++];
			}

			if (header == (byte) 1) {
				int size = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					size |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				if (size > Colfer.colferSizeMax)
					throw new SecurityException(format("colfer: field testdata/bench.Colfer.host size %d exceeds %d UTF-8 bytes", size, Colfer.colferSizeMax));

				int start = i;
				i += size;
				this.host = new String(buf, start, size, StandardCharsets.UTF_8);
				header = buf[i++];
			}

			if (header == (byte) 2) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.port = x;
				header = buf[i++];
			} else if (header == (byte) (2 | 0x80)) {
				int x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					x |= (b & 0x7f) << shift;
					if (shift == 28 || b >= 0) break;
				}
				this.port = -x;
				header = buf[i++];
			}

			if (header == (byte) 3) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.size = x;
				header = buf[i++];
			} else if (header == (byte) (3 | 0x80)) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.size = -x;
				header = buf[i++];
			}

			if (header == (byte) 4) {
				long x = 0;
				for (int shift = 0; true; shift += 7) {
					byte b = buf[i++];
					if (shift == 56 || b >= 0) {
						x |= (b & 0xffL) << shift;
						break;
					}
					x |= (b & 0x7fL) << shift;
				}
				this.hash = x;
				header = buf[i++];
			} else if (header == (byte) (4 | 0x80)) {
				this.hash = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				header = buf[i++];
			}

			if (header == (byte) 5) {
				long x = (buf[i++] & 0xffL) << 56 | (buf[i++] & 0xffL) << 48 | (buf[i++] & 0xffL) << 40 | (buf[i++] & 0xffL) << 32
					| (buf[i++] & 0xffL) << 24 | (buf[i++] & 0xffL) << 16 | (buf[i++] & 0xffL) << 8 | (buf[i++] & 0xffL);
				this.ratio = Double.longBitsToDouble(x);
				header = buf[i++];
			}

			if (header == (byte) 6) {
				this.route = true;
				header = buf[i++];
			}

			if (header != (byte) 0x7f)
				throw new InputMismatchException(format("colfer: unknown header at byte %d", i - 1));
		} finally {
			if (i > end && end - offset < Colfer.colferSizeMax) throw new BufferUnderflowException();
			if (i - offset > Colfer.colferSizeMax)
				throw new SecurityException(format("colfer: serial exceeds %d bytes", Colfer.colferSizeMax));
			if (i > end) throw new BufferUnderflowException();
		}

		return i;
	}

	public long getKey() {
		return this.key;
	}

	public void setKey(long value) {
		this.key = value;
	}

	public String getHost() {
		return this.host;
	}

	public void setHost(String value) {
		this.host = value;
	}

	public int getPort() {
		return this.port;
	}

	public void setPort(int value) {
		this.port = value;
	}

	public long getSize() {
		return this.size;
	}

	public void setSize(long value) {
		this.size = value;
	}

	public long getHash() {
		return this.hash;
	}

	public void setHash(long value) {
		this.hash = value;
	}

	public double getRatio() {
		return this.ratio;
	}

	public void setRatio(double value) {
		this.ratio = value;
	}

	public boolean getRoute() {
		return this.route;
	}

	public void setRoute(boolean value) {
		this.route = value;
	}

	@Override
	public final int hashCode() {
		int h = 1;
		h = 31 * h + (int)(this.key ^ this.key >>> 32);
		if (this.host != null) h = 31 * h + this.host.hashCode();
		h = 31 * h + this.port;
		h = 31 * h + (int)(this.size ^ this.size >>> 32);
		h = 31 * h + (int)(this.hash ^ this.hash >>> 32);
		long _ratioBits = Double.doubleToLongBits(this.ratio);
		h = 31 * h + (int) (_ratioBits ^ _ratioBits >>> 32);
		h = 31 * h + (this.route ? 1231 : 1237);
		return h;
	}

	@Override
	public final boolean equals(Object o) {
		return o instanceof Colfer && equals((Colfer) o);
	}

	public final boolean equals(Colfer o) {
		return o != null && o.getClass() == Colfer.class
			&& this.key == o.key
			&& this.host == null ? o.host == null : this.host.equals(o.host)
			&& this.port == o.port
			&& this.size == o.size
			&& this.hash == o.hash
			&& (this.ratio == o.ratio || (this.ratio != this.ratio && o.ratio != o.ratio))
			&& this.route == o.route;
	}

}
