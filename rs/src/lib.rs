use std::convert::TryInto;

struct State {
    v0: u64,
    v1: u64,
}

impl State {
    fn round(&mut self) {
        self.v0 = self.v0.wrapping_add(self.v1);

        self.v1 = self.v1.rotate_left(13);
        self.v1 ^= self.v0;

        self.v0 = self.v0.rotate_left(35);
        self.v0 = self.v0.wrapping_add(self.v1);

        self.v1 = self.v1.rotate_left(17);
        self.v1 ^= self.v0;

        self.v0 = self.v0.rotate_left(21);
    }
}

pub fn hash(k0: u64, k1: u64, data: &[u8]) -> u64 {
    let bytes = &data;

    let mut s = State {
        v0: k0 ^ 0x736f6d6570736575,
        v1: k1 ^ 0x646f72616e646f6d,
    };

    let mut b = (bytes.len() << 56) as u64;

    let it = bytes.chunks_exact(8);
    let rem = it.remainder();

    for ch in it {
        let m = u64::from_le_bytes(ch[..8].try_into().unwrap());
        s.v1 ^= m;
        s.round();
        s.v0 ^= m;
    }

    match rem.len() {
        0 => (),
        1 => b |= rem[0] as u64,
        2 => b |= u16::from_le_bytes(rem[..2].try_into().unwrap()) as u64,
        3 => b |= u16::from_le_bytes(rem[..2].try_into().unwrap()) as u64 | (rem[2] as u64) << 16,
        4 => b |= u32::from_le_bytes(rem[..4].try_into().unwrap()) as u64,
        5 => b |= u32::from_le_bytes(rem[..4].try_into().unwrap()) as u64 | (rem[4] as u64) << 32,
        6 => {
            b |= u32::from_le_bytes(rem[..4].try_into().unwrap()) as u64
                | (u16::from_le_bytes(rem[4..6].try_into().unwrap()) as u64) << 32
        }
        7 => {
            b |= u32::from_le_bytes(rem[..4].try_into().unwrap()) as u64
                | (u16::from_le_bytes(rem[4..6].try_into().unwrap()) as u64) << 32
                | (rem[6] as u64) << 48
        }
        _ => panic!("len > 7"),
    }

    // last block
    s.v1 ^= b;
    s.round();
    s.v0 ^= b;

    // finalization
    s.v1 ^= 0xff;
    s.round();
    s.v1 = s.v1.rotate_left(32);
    s.round();
    s.v1 = s.v1.rotate_left(32);

    s.v0 ^ s.v1
}

#[cfg(test)]
mod tests {

    use super::*;

    use std::fs::File;
    use std::io::{self, BufRead};
    use std::path::Path;

    #[test]
    fn smoke() {
        let k0: u64 = 0x0706050403020100;
        let k1: u64 = 0x0f0e0d0c0b0a0908;

        let mut b = Vec::new();

        if let Ok(lines) = read_lines("../go/testdata/tsip.txt") {
            for (i, line) in lines.enumerate() {
                if let Ok(want) = line {
                    let want64 = u64::from_str_radix(&want, 16).unwrap();
                    let h = hash(k0, k1, &b);
                    assert_eq!(h, want64);
                    b.push(i as u8);
                }
            }
        }
    }

    fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
    where
        P: AsRef<Path>,
    {
        let file = File::open(filename)?;
        Ok(io::BufReader::new(file).lines())
    }
}
