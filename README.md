# ceph-crush-ana

## Imitate ceph's crush algorithm in golang

# ceph build rpm
wget -P ~/rpmbuild/SOURCES/ http://download.ceph.com/tarballs/ceph-11.2.0.tar.gz

tar --strip-components=1 -C ~/rpmbuild/SPECS/ --no-anchored -zxvf ~/rpmbuild/SOURCES/ceph-11.2.0.tar.gz "ceph.spec"

yum install rpm-build rpmdevtools -y

yum install java-devel sharutils checkpolicy selinux-policy-devel /usr/share/selinux/devel/policyhelp   boost-devel boost-python cmake cryptsetup fuse-devel gcc-c++ gperftools-devel jq leveldb-devel libaio-devel libatomic_ops-devel libblkid-devel libcurl-devel libudev-devel libtool libxml2-devel python-devel python-nose python-requests python-sphinx python-virtualenv snappy-devel valgrind-devel xfsprogs-devel xmlstarlet yasm boost-random nss-devel keyutils-libs-devel openldap-devel openssl-devel redhat-lsb-core Cython python34-devel python34-setuptools python34-Cython lttng-ust-devel -y

chmown root:root ~/rpmbuild/SPECS/ceph.spec
rpmbuild -ba ~/rpmbuild/SPECS/ceph.spec

# build ceph from source

wget https://download.ceph.com/tarballs/ceph_10.2.9.orig.tar.gz

tar -xvf ceph_10.2.9.orig.tar.gz

yum install java-devel sharutils checkpolicy selinux-policy-devel /usr/share/selinux/devel/policyhelp   boost-devel boost-python cmake cryptsetup fuse-devel gcc-c++ gperftools-devel jq leveldb-devel libaio-devel libatomic_ops-devel libblkid-devel libcurl-devel libudev-devel libtool libxml2-devel python-devel python-nose python-requests python-sphinx python-virtualenv snappy-devel valgrind-devel xfsprogs-devel xmlstarlet yasm boost-random nss-devel keyutils-libs-devel openldap-devel openssl-devel redhat-lsb-core Cython python34-devel python34-setuptools python34-Cython lttng-ust-devel -y

cd ceph_10.2.9

./autogen.sh

./configure

make -j5

bug:

1. mds/MDCache.CC +9130

bool was_replay = mdr->client_request && mdr->client_request->is_replay();

was_replay 定义但未使用

2. rgw/rgw_rest_swift.cc:8:22: fatal error: ceph_ver.h: No such file or directory
 
 #include "ceph_ver.h"

3. rgw/rgw_swift_auth.cc:144:15: warning: unused variable 'token_tag' [-Wunused-variable]
   const char *token_tag = "rgwtk";
   
   
   
## install ceph

ceph-deploy uninstall ceph0 ceph1 ceph2

ceph-deploy purge ceph0 ceph1 ceph2

ceph-deploy purgedata ceph0 ceph1 ceph2

ceph-deploy forgetkeys

ceph-deploy install --repo-url http://mirrors.163.com/ceph/rpm-jewel/el7 --nogpgcheck --gpg-url http://mirrors.163.com/ceph/keys/release.asc ceph0 ceph1 ceph2

ceph-deploy --overwrite-conf mon create-initial

vim cleanOsd.sh
```
#!/bin/sh

hosts="ceph0 ceph1 ceph2"
dev="b c"

for hostn in $hosts
do
        for i in $dev
        do
                ceph-deploy disk zap ${hostn}:sd${i}
        done
done
```
vim createOsd.sh
```
#!/bin/sh

hosts="ceph0 ceph1 ceph2"
dev="b c"

for hostn in $hosts
do
        for i in $dev
        do
                ceph-deploy osd create ${hostn}:sd${i}
        done
done
```

ceph-deploy admin ceph0 ceph1 ceph2

sudo chmod +r /etc/ceph/ceph.client.admin.keyring

ceph -s

ceph osd pool create fiotest 1024 1024
