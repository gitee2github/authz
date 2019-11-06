#Global macro or variable
#%define _unpackaged_files_terminate_build 0

#Basic Information
Name:           authz
Version:        0.1
Release:        16
Summary:        a isula auth plugin for RBAC
License:        Mulan PSL v1
Source0:        authz-plugin.tar.gz
BuildRoot:      %{_tmppath}/authz-root

#Dependency
BuildRequires: golang >= 1.8
BuildRequires: glibc-static

%description
Work with isulad daemon that enables TLS. It brings the support of RBAC.

#Build sections
%prep
export RPM_BUILD_SOURCE=%_topdir/SOURCES
export RPM_BUILD_DIR=%_topdir/BUILD

cd $RPM_BUILD_DIR
mkdir -p $RPM_BUILD_DIR/src/isula.org/authz && cd $RPM_BUILD_DIR/src/isula.org/authz
gzip -dc $RPM_BUILD_SOURCE/authz-plugin.tar.gz | tar -xvvf -

%build
cd $RPM_BUILD_DIR/src/isula.org/authz
export GOPATH=%_topdir/BUILD
make

%install
cd $RPM_BUILD_DIR/src/isula.org/authz
mkdir -p $RPM_BUILD_ROOT/usr/lib/isulad/
mkdir -p $RPM_BUILD_ROOT/lib/systemd/system/
mkdir -p $RPM_BUILD_ROOT/var/lib/authz-broker/

cp bin/authz-broker $RPM_BUILD_ROOT/usr/lib/isulad/
cp systemd/authz.service $RPM_BUILD_ROOT/lib/systemd/system/

chmod 0750 $RPM_BUILD_ROOT/usr/lib/isulad/authz-broker

#Install and uninstall scripts
%pre

%preun
%systemd_preun authz

%post
if [ ! -d "/var/lib/authz-broker" ]; then
	mkdir -p /var/lib/authz-broker
fi
chmod 0750 /var/lib/authz-broker
if [ ! -f "/var/lib/authz-broker/policy.json" ]; then
	cat > /var/lib/authz-broker/policy.json << EOF
{"name":"policy_root","users":[""],"actions":[""]}
EOF
fi
chmod 0640 /var/lib/authz-broker/policy.json

%postun

#Files list
%files
%attr(550,root,root) /usr/lib/isulad
%attr(550,root,root) /usr/lib/isulad/authz-broker
%attr(640,root,root) /lib/systemd/system/authz.service

#Clean section
%clean
rm -rfv %{buildroot}

%changelog
* Tue Oct 23 2018 Zhangsong<zhangsong34@huawei.com> - 0.1.0-6
- Type:enhancement
- ID:NA
- SUG:restart
- DESC:support isulad update permission control for authz

* Fri Aug 10 2018 Liruilin<liruilin4@huawei.com> - 0.1.0-5
- Type:enhancement
- ID:NA
- SUG:restart
- DESC:add --pidfile option

* Mon May 21 2018 Liruilin<liruilin4@huawei.com> - 0.1.0-4
- Type:enhancement
- ID:NA
- SUG:restart
- DESC:create config file dir in %post

* Mon May 7 2018 Liruilin<liruilin4@huawei.com> - 0.1.0-3
- Type:enhancement
- ID:NA
- SUG:restart
- DESC:move policy.json to %post

* Mon May 7 2018 Liruilin<liruilin4@huawei.com> - 0.1.0-2
- Type:enhancement
- ID:NA
- SUG:restart
- DESC:repair systemd service files

* Sat Apr 21 2018 Caoruidong<caoruidong@huawei.com> - 0.1.0-1
- Type:new-packages
- ID:NA
- SUG:restart
- DESC:add a new auth plugin for isulad
