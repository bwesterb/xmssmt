XMSSMT commandline tool
=======================

This is a commandline tool to sign and verify messages using the
post-quantum stateful hash-based signature-scheme XMSSMT described in
[rfc8391](https://tools.ietf.org/html/rfc8391).

Installing
----------

To install `xmssmt`, get [Go](https://golang.org/) and run

    GO111MODULE=on go get github.com/bwesterb/xmssmt

Usage
-----

#### Generating a keypair

To generate an XMSSMT keypair, run

    xmssmt generate

This will generate a random `XMSSMT-SHAKE_40/4_256` keypair.  The secret key is
stored in two files: `xmssmt.key` and `xmssmt.key.cache`.  You should keep
both files secret, never copy them and never restore them from a backup.
(See below.) The public key is stores in `xmssmt.pub`.

You can specify a different instance of XMSSMT with `-a`, for instance

    xmssmt generate -a XMSSMT-SHA2_20/2_512

Run `xmssmt algs` to list the named instances.  (See below for the considerations
when choosing an XMSSMT instance.)

#### Signing

To create an XMSSMT signature on `some-file`, run

    xmssmt sign -f some-file

This will create an XMSSMT signature `some-file.xmssmt-signature` using the
secret key `xmssmt.key`.

A different secret key and signature output file can be specified with flags:

    xmssmt sign -f some-file -s path/to/secret-key -o path/to/write/signature/to

#### Verifying

To verify the XMSSMT signature `some-file.xmsssmt-signature` on `some-file`, run

    xmssmt verify -f some-file

It will look for the public key in the file `xmssmt.pub`.

With flags one can specify the files to read the signature and public key from.  Eg:

    xmssmt verify -f some-file -S the-signature -p path-to-public-key

Considerations
--------------

### State

XMSSMT (in contrast to its sibling [SPHINCS+](https://sphincs.org/)) is stateful:
every signature has a sequence number and a sequence number
[should](https://eprint.iacr.org/2016/1042.pdf) not be reused.
There is also a maximum signature sequence number (dependant on the exact
XMSSMT instance).
The first free signature sequence number is stored in the secret key
file `xmssmt.key`, which is incremented on every signature issued.  Thus

 * You should **not copy** the secret key file, for otherwise signature
   sequence numbers might be reused.
 * You should **never restore** a secret key file from a backup, for again,
   otherwise signature sequence numbers might be reused.

### Cache

Without cache, creating a XMSSMT signature is about as expensive as generating
a keypair.  Almost all computations between two signatures (which are close in
sequence number) can be reused.  To this end `xmssmt` stores these values
in the `.key.cache` file.  This makes creating a signature even significantly
faster than verifying one (if cached).

With the default XMSSMT instance, signatures are cached in batches
("subtrees") of 1024.  So, the first 1024 signatures are quick to create.
The 1025th signature takes (with the default instance) a fourth of the
key generation time and the next 1023 signatures are again very fast to create.

### Instance & parameters

An XMSSMT instance has five main parameters

 * The **hash function** used.  Either SHAKE or SHA2.  The XMSSMT authors prefer
   SHAKE and it's significantly faster than SHA2 (except for n=512.)
 * **n** is the main security parameter and is either 128, 256 or 512 bits.
   512 bit signatures are at least twice as large (see `xmssmt algs`),
   and are approximately tree times as slow to create and verify.
   If you're unsure, use 256 bit.  Use 512 bit if
   
    1. you want your signatures to be trustworthy for at least a 100 years *and*
    2. you believe that performance per watt will keep increasing exponentially.

 * **tree height** determines the maximum number of signatures that can be created
   with a keypair.  With tree height t, one can create 2^t signatures.  So with the
   default tree height 40, we have about a trillion signatures.  With the other
   parameters fixed, a higher tree height will exponentially increase key
   generation time and secret key (cache) size.  In contrast, signature size
   and signing/verifying-times will not change by much.
 * **d** is, in effect, a trade-off between

     1. signature size and signature verification times and
     2. secret key (cache) size and keypair generation time.
   
   If unsure, pick `d = (tree height) / 10`. Then the keypair will only take a
   few seconds to generate; the secret key (cache) is less than a megabyte;
   signatures are cached in batches of 1024 and still have an acceptable
   size (for the tree height).
   
   If long key generation time, ~250MB secret key cache and slow signing
   every millionth of a signature is not a problem, consider `d = (tree height) / 20`.
 * **w** is a trade-off between

    1. signature size
    2. signing/verification/key generation time.

   The default is `w=16`.  The RFC only lists instances with `w=16`.

The parameters (for `w=16`) are formatted as follows in the name

    XMSSMT-(hash func)_(tree height)/(d)_(n in bits)

The special case `d=1` is formatted as

    XMSS-(hash func)_(tree height)_(n in bits)

The parameter `w` can be specified by suffixing `_w(value of w)`, i.e.:

    XMSSMT-SHAKE_60/12_512_w256

See also
--------

 * [atumd](https://github.com/bwesterb/atumd), a timestamping server that uses XMSSMT
 * [go-xmssmt](https://github.com/bwesterb/go-xmssmt), a Go package that implements XMSSMT.

