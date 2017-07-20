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

tar -xvf ceph-0.94.3.tar.gz

yum install java-devel sharutils checkpolicy selinux-policy-devel /usr/share/selinux/devel/policyhelp   boost-devel boost-python cmake cryptsetup fuse-devel gcc-c++ gperftools-devel jq leveldb-devel libaio-devel libatomic_ops-devel libblkid-devel libcurl-devel libudev-devel libtool libxml2-devel python-devel python-nose python-requests python-sphinx python-virtualenv snappy-devel valgrind-devel xfsprogs-devel xmlstarlet yasm boost-random nss-devel keyutils-libs-devel openldap-devel openssl-devel redhat-lsb-core Cython python34-devel python34-setuptools python34-Cython lttng-ust-devel -y

cd ceph_10.2.9

./autogen.sh

./configure

make -j5

bug:
1. mds/MDCache.CC +9130
bool was_replay = mdr->client_request && mdr->client_request->is_replay();
was_replay 定义但未使用
