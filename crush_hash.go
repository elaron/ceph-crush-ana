package main

// CrushHashSeed is the hash_seed
const CrushHashSeed uint32 = 1315423911

/*
 * Robert Jenkins' function for mixing 32-bit values
 * http://burtleburtle.net/bob/hash/evahash.html
 * a, b = random bits, c = input and output
 */
func crushHashmix(a, b, c *uint32) {
	(*a) -= *b
	(*a) -= *c
	*a = *a ^ (*c >> 13)
	*b -= *c
	*b -= *a
	*b = *b ^ (*a << 8)
	*c -= *a
	*c -= *b
	*c = *c ^ (*b >> 13)
	*a -= *b
	*a -= *c
	*a = *a ^ (*c >> 12)
	*b -= *c
	*b -= *a
	*b = *b ^ (*a << 16)
	*c -= *a
	*c -= *b
	*c = *c ^ (*b >> 5)
	*a -= *b
	*a -= *c
	*a = *a ^ (*c >> 3)
	*b -= *c
	*b -= *a
	*b = *b ^ (*a << 10)
	*c -= *a
	*c -= *b
	*c = *c ^ (*b >> 15)
}

//crushHash32_3 represents hash with 3 parameters
func crushHash32_3(a, b, c uint32) uint32 {
	var hash, x, y uint32
	hash = CrushHashSeed ^ a ^ b ^ c
	x = 231232
	y = 1232
	crushHashmix(&a, &b, &hash)
	crushHashmix(&c, &x, &hash)
	crushHashmix(&y, &a, &hash)
	crushHashmix(&b, &x, &hash)
	crushHashmix(&y, &c, &hash)
	return hash
}
