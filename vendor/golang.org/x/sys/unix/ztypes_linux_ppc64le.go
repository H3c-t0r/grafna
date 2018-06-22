// cgo -godefs -- -Wall -Werror -static -I/tmp/include linux/types.go | go run mkpost.go
// Code generated by the command above; see README.md. DO NOT EDIT.

// +build ppc64le,linux

package unix

const (
	sizeofPtr      = 0x8
	sizeofShort    = 0x2
	sizeofInt      = 0x4
	sizeofLong     = 0x8
	sizeofLongLong = 0x8
	PathMax        = 0x1000
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

type Timespec struct {
	Sec  int64
	Nsec int64
}

type Timeval struct {
	Sec  int64
	Usec int64
}

type Timex struct {
	Modes     uint32
	_         [4]byte
	Offset    int64
	Freq      int64
	Maxerror  int64
	Esterror  int64
	Status    int32
	_         [4]byte
	Constant  int64
	Precision int64
	Tolerance int64
	Time      Timeval
	Tick      int64
	Ppsfreq   int64
	Jitter    int64
	Shift     int32
	_         [4]byte
	Stabil    int64
	Jitcnt    int64
	Calcnt    int64
	Errcnt    int64
	Stbcnt    int64
	Tai       int32
	_         [44]byte
}

type Time_t int64

type Tms struct {
	Utime  int64
	Stime  int64
	Cutime int64
	Cstime int64
}

type Utimbuf struct {
	Actime  int64
	Modtime int64
}

type Rusage struct {
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
}

type Rlimit struct {
	Cur uint64
	Max uint64
}

type _Gid_t uint32

type Stat_t struct {
	Dev     uint64
	Ino     uint64
	Nlink   uint64
	Mode    uint32
	Uid     uint32
	Gid     uint32
	_       int32
	Rdev    uint64
	Size    int64
	Blksize int64
	Blocks  int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	_       uint64
	_       uint64
	_       uint64
}

type StatxTimestamp struct {
	Sec  int64
	Nsec uint32
	_    int32
}

type Statx_t struct {
	Mask            uint32
	Blksize         uint32
	Attributes      uint64
	Nlink           uint32
	Uid             uint32
	Gid             uint32
	Mode            uint16
	_               [1]uint16
	Ino             uint64
	Size            uint64
	Blocks          uint64
	Attributes_mask uint64
	Atime           StatxTimestamp
	Btime           StatxTimestamp
	Ctime           StatxTimestamp
	Mtime           StatxTimestamp
	Rdev_major      uint32
	Rdev_minor      uint32
	Dev_major       uint32
	Dev_minor       uint32
	_               [14]uint64
}

type Dirent struct {
	Ino    uint64
	Off    int64
	Reclen uint16
	Type   uint8
	Name   [256]uint8
	_      [5]byte
}

type Fsid struct {
	Val [2]int32
}

type Flock_t struct {
	Type   int16
	Whence int16
	_      [4]byte
	Start  int64
	Len    int64
	Pid    int32
	_      [4]byte
}

type FscryptPolicy struct {
	Version                   uint8
	Contents_encryption_mode  uint8
	Filenames_encryption_mode uint8
	Flags                     uint8
	Master_key_descriptor     [8]uint8
}

type FscryptKey struct {
	Mode uint32
	Raw  [64]uint8
	Size uint32
}

type KeyctlDHParams struct {
	Private int32
	Prime   int32
	Base    int32
}

const (
	FADV_NORMAL     = 0x0
	FADV_RANDOM     = 0x1
	FADV_SEQUENTIAL = 0x2
	FADV_WILLNEED   = 0x3
	FADV_DONTNEED   = 0x4
	FADV_NOREUSE    = 0x5
)

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
}

type RawSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
}

type RawSockaddrUnix struct {
	Family uint16
	Path   [108]int8
}

type RawSockaddrLinklayer struct {
	Family   uint16
	Protocol uint16
	Ifindex  int32
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]uint8
}

type RawSockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type RawSockaddrHCI struct {
	Family  uint16
	Dev     uint16
	Channel uint16
}

type RawSockaddrL2 struct {
	Family      uint16
	Psm         uint16
	Bdaddr      [6]uint8
	Cid         uint16
	Bdaddr_type uint8
	_           [1]byte
}

type RawSockaddrCAN struct {
	Family  uint16
	_       [2]byte
	Ifindex int32
	Addr    [8]byte
}

type RawSockaddrALG struct {
	Family uint16
	Type   [14]uint8
	Feat   uint32
	Mask   uint32
	Name   [64]uint8
}

type RawSockaddrVM struct {
	Family    uint16
	Reserved1 uint16
	Port      uint32
	Cid       uint32
	Zero      [4]uint8
}

type RawSockaddr struct {
	Family uint16
	Data   [14]uint8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [96]uint8
}

type _Socklen uint32

type Linger struct {
	Onoff  int32
	Linger int32
}

type Iovec struct {
	Base *byte
	Len  uint64
}

type IPMreq struct {
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
}

type IPMreqn struct {
	Multiaddr [4]byte /* in_addr */
	Address   [4]byte /* in_addr */
	Ifindex   int32
}

type IPv6Mreq struct {
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
}

type PacketMreq struct {
	Ifindex int32
	Type    uint16
	Alen    uint16
	Address [8]uint8
}

type Msghdr struct {
	Name       *byte
	Namelen    uint32
	_          [4]byte
	Iov        *Iovec
	Iovlen     uint64
	Control    *byte
	Controllen uint64
	Flags      int32
	_          [4]byte
}

type Cmsghdr struct {
	Len   uint64
	Level int32
	Type  int32
}

type Inet4Pktinfo struct {
	Ifindex  int32
	Spec_dst [4]byte /* in_addr */
	Addr     [4]byte /* in_addr */
}

type Inet6Pktinfo struct {
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
}

type IPv6MTUInfo struct {
	Addr RawSockaddrInet6
	Mtu  uint32
}

type ICMPv6Filter struct {
	Data [8]uint32
}

type Ucred struct {
	Pid int32
	Uid uint32
	Gid uint32
}

type TCPInfo struct {
	State          uint8
	Ca_state       uint8
	Retransmits    uint8
	Probes         uint8
	Backoff        uint8
	Options        uint8
	_              [2]byte
	Rto            uint32
	Ato            uint32
	Snd_mss        uint32
	Rcv_mss        uint32
	Unacked        uint32
	Sacked         uint32
	Lost           uint32
	Retrans        uint32
	Fackets        uint32
	Last_data_sent uint32
	Last_ack_sent  uint32
	Last_data_recv uint32
	Last_ack_recv  uint32
	Pmtu           uint32
	Rcv_ssthresh   uint32
	Rtt            uint32
	Rttvar         uint32
	Snd_ssthresh   uint32
	Snd_cwnd       uint32
	Advmss         uint32
	Reordering     uint32
	Rcv_rtt        uint32
	Rcv_space      uint32
	Total_retrans  uint32
}

const (
	SizeofSockaddrInet4     = 0x10
	SizeofSockaddrInet6     = 0x1c
	SizeofSockaddrAny       = 0x70
	SizeofSockaddrUnix      = 0x6e
	SizeofSockaddrLinklayer = 0x14
	SizeofSockaddrNetlink   = 0xc
	SizeofSockaddrHCI       = 0x6
	SizeofSockaddrL2        = 0xe
	SizeofSockaddrCAN       = 0x10
	SizeofSockaddrALG       = 0x58
	SizeofSockaddrVM        = 0x10
	SizeofLinger            = 0x8
	SizeofIovec             = 0x10
	SizeofIPMreq            = 0x8
	SizeofIPMreqn           = 0xc
	SizeofIPv6Mreq          = 0x14
	SizeofPacketMreq        = 0x10
	SizeofMsghdr            = 0x38
	SizeofCmsghdr           = 0x10
	SizeofInet4Pktinfo      = 0xc
	SizeofInet6Pktinfo      = 0x14
	SizeofIPv6MTUInfo       = 0x20
	SizeofICMPv6Filter      = 0x20
	SizeofUcred             = 0xc
	SizeofTCPInfo           = 0x68
)

const (
	IFA_UNSPEC           = 0x0
	IFA_ADDRESS          = 0x1
	IFA_LOCAL            = 0x2
	IFA_LABEL            = 0x3
	IFA_BROADCAST        = 0x4
	IFA_ANYCAST          = 0x5
	IFA_CACHEINFO        = 0x6
	IFA_MULTICAST        = 0x7
	IFLA_UNSPEC          = 0x0
	IFLA_ADDRESS         = 0x1
	IFLA_BROADCAST       = 0x2
	IFLA_IFNAME          = 0x3
	IFLA_MTU             = 0x4
	IFLA_LINK            = 0x5
	IFLA_QDISC           = 0x6
	IFLA_STATS           = 0x7
	IFLA_COST            = 0x8
	IFLA_PRIORITY        = 0x9
	IFLA_MASTER          = 0xa
	IFLA_WIRELESS        = 0xb
	IFLA_PROTINFO        = 0xc
	IFLA_TXQLEN          = 0xd
	IFLA_MAP             = 0xe
	IFLA_WEIGHT          = 0xf
	IFLA_OPERSTATE       = 0x10
	IFLA_LINKMODE        = 0x11
	IFLA_LINKINFO        = 0x12
	IFLA_NET_NS_PID      = 0x13
	IFLA_IFALIAS         = 0x14
	IFLA_NUM_VF          = 0x15
	IFLA_VFINFO_LIST     = 0x16
	IFLA_STATS64         = 0x17
	IFLA_VF_PORTS        = 0x18
	IFLA_PORT_SELF       = 0x19
	IFLA_AF_SPEC         = 0x1a
	IFLA_GROUP           = 0x1b
	IFLA_NET_NS_FD       = 0x1c
	IFLA_EXT_MASK        = 0x1d
	IFLA_PROMISCUITY     = 0x1e
	IFLA_NUM_TX_QUEUES   = 0x1f
	IFLA_NUM_RX_QUEUES   = 0x20
	IFLA_CARRIER         = 0x21
	IFLA_PHYS_PORT_ID    = 0x22
	IFLA_CARRIER_CHANGES = 0x23
	IFLA_PHYS_SWITCH_ID  = 0x24
	IFLA_LINK_NETNSID    = 0x25
	IFLA_PHYS_PORT_NAME  = 0x26
	IFLA_PROTO_DOWN      = 0x27
	IFLA_GSO_MAX_SEGS    = 0x28
	IFLA_GSO_MAX_SIZE    = 0x29
	IFLA_PAD             = 0x2a
	IFLA_XDP             = 0x2b
	IFLA_EVENT           = 0x2c
	IFLA_NEW_NETNSID     = 0x2d
	IFLA_IF_NETNSID      = 0x2e
	IFLA_MAX             = 0x31
	RT_SCOPE_UNIVERSE    = 0x0
	RT_SCOPE_SITE        = 0xc8
	RT_SCOPE_LINK        = 0xfd
	RT_SCOPE_HOST        = 0xfe
	RT_SCOPE_NOWHERE     = 0xff
	RT_TABLE_UNSPEC      = 0x0
	RT_TABLE_COMPAT      = 0xfc
	RT_TABLE_DEFAULT     = 0xfd
	RT_TABLE_MAIN        = 0xfe
	RT_TABLE_LOCAL       = 0xff
	RT_TABLE_MAX         = 0xffffffff
	RTA_UNSPEC           = 0x0
	RTA_DST              = 0x1
	RTA_SRC              = 0x2
	RTA_IIF              = 0x3
	RTA_OIF              = 0x4
	RTA_GATEWAY          = 0x5
	RTA_PRIORITY         = 0x6
	RTA_PREFSRC          = 0x7
	RTA_METRICS          = 0x8
	RTA_MULTIPATH        = 0x9
	RTA_FLOW             = 0xb
	RTA_CACHEINFO        = 0xc
	RTA_TABLE            = 0xf
	RTN_UNSPEC           = 0x0
	RTN_UNICAST          = 0x1
	RTN_LOCAL            = 0x2
	RTN_BROADCAST        = 0x3
	RTN_ANYCAST          = 0x4
	RTN_MULTICAST        = 0x5
	RTN_BLACKHOLE        = 0x6
	RTN_UNREACHABLE      = 0x7
	RTN_PROHIBIT         = 0x8
	RTN_THROW            = 0x9
	RTN_NAT              = 0xa
	RTN_XRESOLVE         = 0xb
	RTNLGRP_NONE         = 0x0
	RTNLGRP_LINK         = 0x1
	RTNLGRP_NOTIFY       = 0x2
	RTNLGRP_NEIGH        = 0x3
	RTNLGRP_TC           = 0x4
	RTNLGRP_IPV4_IFADDR  = 0x5
	RTNLGRP_IPV4_MROUTE  = 0x6
	RTNLGRP_IPV4_ROUTE   = 0x7
	RTNLGRP_IPV4_RULE    = 0x8
	RTNLGRP_IPV6_IFADDR  = 0x9
	RTNLGRP_IPV6_MROUTE  = 0xa
	RTNLGRP_IPV6_ROUTE   = 0xb
	RTNLGRP_IPV6_IFINFO  = 0xc
	RTNLGRP_IPV6_PREFIX  = 0x12
	RTNLGRP_IPV6_RULE    = 0x13
	RTNLGRP_ND_USEROPT   = 0x14
	SizeofNlMsghdr       = 0x10
	SizeofNlMsgerr       = 0x14
	SizeofRtGenmsg       = 0x1
	SizeofNlAttr         = 0x4
	SizeofRtAttr         = 0x4
	SizeofIfInfomsg      = 0x10
	SizeofIfAddrmsg      = 0x8
	SizeofRtMsg          = 0xc
	SizeofRtNexthop      = 0x8
)

type NlMsghdr struct {
	Len   uint32
	Type  uint16
	Flags uint16
	Seq   uint32
	Pid   uint32
}

type NlMsgerr struct {
	Error int32
	Msg   NlMsghdr
}

type RtGenmsg struct {
	Family uint8
}

type NlAttr struct {
	Len  uint16
	Type uint16
}

type RtAttr struct {
	Len  uint16
	Type uint16
}

type IfInfomsg struct {
	Family uint8
	_      uint8
	Type   uint16
	Index  int32
	Flags  uint32
	Change uint32
}

type IfAddrmsg struct {
	Family    uint8
	Prefixlen uint8
	Flags     uint8
	Scope     uint8
	Index     uint32
}

type RtMsg struct {
	Family   uint8
	Dst_len  uint8
	Src_len  uint8
	Tos      uint8
	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8
	Flags    uint32
}

type RtNexthop struct {
	Len     uint16
	Flags   uint8
	Hops    uint8
	Ifindex int32
}

const (
	SizeofSockFilter = 0x8
	SizeofSockFprog  = 0x10
)

type SockFilter struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

type SockFprog struct {
	Len    uint16
	_      [6]byte
	Filter *SockFilter
}

type InotifyEvent struct {
	Wd     int32
	Mask   uint32
	Cookie uint32
	Len    uint32
}

const SizeofInotifyEvent = 0x10

type PtraceRegs struct {
	Gpr       [32]uint64
	Nip       uint64
	Msr       uint64
	Orig_gpr3 uint64
	Ctr       uint64
	Link      uint64
	Xer       uint64
	Ccr       uint64
	Softe     uint64
	Trap      uint64
	Dar       uint64
	Dsisr     uint64
	Result    uint64
}

type FdSet struct {
	Bits [16]int64
}

type Sysinfo_t struct {
	Uptime    int64
	Loads     [3]uint64
	Totalram  uint64
	Freeram   uint64
	Sharedram uint64
	Bufferram uint64
	Totalswap uint64
	Freeswap  uint64
	Procs     uint16
	Pad       uint16
	_         [4]byte
	Totalhigh uint64
	Freehigh  uint64
	Unit      uint32
	_         [0]uint8
	_         [4]byte
}

type Utsname struct {
	Sysname    [65]byte
	Nodename   [65]byte
	Release    [65]byte
	Version    [65]byte
	Machine    [65]byte
	Domainname [65]byte
}

type Ustat_t struct {
	Tfree  int32
	_      [4]byte
	Tinode uint64
	Fname  [6]uint8
	Fpack  [6]uint8
	_      [4]byte
}

type EpollEvent struct {
	Events uint32
	_      int32
	Fd     int32
	Pad    int32
}

const (
	AT_EMPTY_PATH   = 0x1000
	AT_FDCWD        = -0x64
	AT_NO_AUTOMOUNT = 0x800
	AT_REMOVEDIR    = 0x200

	AT_STATX_SYNC_AS_STAT = 0x0
	AT_STATX_FORCE_SYNC   = 0x2000
	AT_STATX_DONT_SYNC    = 0x4000

	AT_SYMLINK_FOLLOW   = 0x400
	AT_SYMLINK_NOFOLLOW = 0x100
)

type PollFd struct {
	Fd      int32
	Events  int16
	Revents int16
}

const (
	POLLIN    = 0x1
	POLLPRI   = 0x2
	POLLOUT   = 0x4
	POLLRDHUP = 0x2000
	POLLERR   = 0x8
	POLLHUP   = 0x10
	POLLNVAL  = 0x20
)

type Sigset_t struct {
	Val [16]uint64
}

const RNDGETENTCNT = 0x40045200

const PERF_IOC_FLAG_GROUP = 0x1

type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [19]uint8
	Line   uint8
	Ispeed uint32
	Ospeed uint32
}

type Winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type Taskstats struct {
	Version                   uint16
	_                         [2]byte
	Ac_exitcode               uint32
	Ac_flag                   uint8
	Ac_nice                   uint8
	_                         [6]byte
	Cpu_count                 uint64
	Cpu_delay_total           uint64
	Blkio_count               uint64
	Blkio_delay_total         uint64
	Swapin_count              uint64
	Swapin_delay_total        uint64
	Cpu_run_real_total        uint64
	Cpu_run_virtual_total     uint64
	Ac_comm                   [32]uint8
	Ac_sched                  uint8
	Ac_pad                    [3]uint8
	_                         [4]byte
	Ac_uid                    uint32
	Ac_gid                    uint32
	Ac_pid                    uint32
	Ac_ppid                   uint32
	Ac_btime                  uint32
	_                         [4]byte
	Ac_etime                  uint64
	Ac_utime                  uint64
	Ac_stime                  uint64
	Ac_minflt                 uint64
	Ac_majflt                 uint64
	Coremem                   uint64
	Virtmem                   uint64
	Hiwater_rss               uint64
	Hiwater_vm                uint64
	Read_char                 uint64
	Write_char                uint64
	Read_syscalls             uint64
	Write_syscalls            uint64
	Read_bytes                uint64
	Write_bytes               uint64
	Cancelled_write_bytes     uint64
	Nvcsw                     uint64
	Nivcsw                    uint64
	Ac_utimescaled            uint64
	Ac_stimescaled            uint64
	Cpu_scaled_run_real_total uint64
	Freepages_count           uint64
	Freepages_delay_total     uint64
}

const (
	TASKSTATS_CMD_UNSPEC                  = 0x0
	TASKSTATS_CMD_GET                     = 0x1
	TASKSTATS_CMD_NEW                     = 0x2
	TASKSTATS_TYPE_UNSPEC                 = 0x0
	TASKSTATS_TYPE_PID                    = 0x1
	TASKSTATS_TYPE_TGID                   = 0x2
	TASKSTATS_TYPE_STATS                  = 0x3
	TASKSTATS_TYPE_AGGR_PID               = 0x4
	TASKSTATS_TYPE_AGGR_TGID              = 0x5
	TASKSTATS_TYPE_NULL                   = 0x6
	TASKSTATS_CMD_ATTR_UNSPEC             = 0x0
	TASKSTATS_CMD_ATTR_PID                = 0x1
	TASKSTATS_CMD_ATTR_TGID               = 0x2
	TASKSTATS_CMD_ATTR_REGISTER_CPUMASK   = 0x3
	TASKSTATS_CMD_ATTR_DEREGISTER_CPUMASK = 0x4
)

type CGroupStats struct {
	Sleeping        uint64
	Running         uint64
	Stopped         uint64
	Uninterruptible uint64
	Io_wait         uint64
}

const (
	CGROUPSTATS_CMD_UNSPEC        = 0x3
	CGROUPSTATS_CMD_GET           = 0x4
	CGROUPSTATS_CMD_NEW           = 0x5
	CGROUPSTATS_TYPE_UNSPEC       = 0x0
	CGROUPSTATS_TYPE_CGROUP_STATS = 0x1
	CGROUPSTATS_CMD_ATTR_UNSPEC   = 0x0
	CGROUPSTATS_CMD_ATTR_FD       = 0x1
)

type Genlmsghdr struct {
	Cmd      uint8
	Version  uint8
	Reserved uint16
}

const (
	CTRL_CMD_UNSPEC            = 0x0
	CTRL_CMD_NEWFAMILY         = 0x1
	CTRL_CMD_DELFAMILY         = 0x2
	CTRL_CMD_GETFAMILY         = 0x3
	CTRL_CMD_NEWOPS            = 0x4
	CTRL_CMD_DELOPS            = 0x5
	CTRL_CMD_GETOPS            = 0x6
	CTRL_CMD_NEWMCAST_GRP      = 0x7
	CTRL_CMD_DELMCAST_GRP      = 0x8
	CTRL_CMD_GETMCAST_GRP      = 0x9
	CTRL_ATTR_UNSPEC           = 0x0
	CTRL_ATTR_FAMILY_ID        = 0x1
	CTRL_ATTR_FAMILY_NAME      = 0x2
	CTRL_ATTR_VERSION          = 0x3
	CTRL_ATTR_HDRSIZE          = 0x4
	CTRL_ATTR_MAXATTR          = 0x5
	CTRL_ATTR_OPS              = 0x6
	CTRL_ATTR_MCAST_GROUPS     = 0x7
	CTRL_ATTR_OP_UNSPEC        = 0x0
	CTRL_ATTR_OP_ID            = 0x1
	CTRL_ATTR_OP_FLAGS         = 0x2
	CTRL_ATTR_MCAST_GRP_UNSPEC = 0x0
	CTRL_ATTR_MCAST_GRP_NAME   = 0x1
	CTRL_ATTR_MCAST_GRP_ID     = 0x2
)

type cpuMask uint64

const (
	_CPU_SETSIZE = 0x400
	_NCPUBITS    = 0x40
)

const (
	BDADDR_BREDR     = 0x0
	BDADDR_LE_PUBLIC = 0x1
	BDADDR_LE_RANDOM = 0x2
)

type PerfEventAttr struct {
	Type               uint32
	Size               uint32
	Config             uint64
	Sample             uint64
	Sample_type        uint64
	Read_format        uint64
	Bits               uint64
	Wakeup             uint32
	Bp_type            uint32
	Ext1               uint64
	Ext2               uint64
	Branch_sample_type uint64
	Sample_regs_user   uint64
	Sample_stack_user  uint32
	Clockid            int32
	Sample_regs_intr   uint64
	Aux_watermark      uint32
	_                  uint32
}

type PerfEventMmapPage struct {
	Version        uint32
	Compat_version uint32
	Lock           uint32
	Index          uint32
	Offset         int64
	Time_enabled   uint64
	Time_running   uint64
	Capabilities   uint64
	Pmc_width      uint16
	Time_shift     uint16
	Time_mult      uint32
	Time_offset    uint64
	Time_zero      uint64
	Size           uint32
	_              [948]uint8
	Data_head      uint64
	Data_tail      uint64
	Data_offset    uint64
	Data_size      uint64
	Aux_head       uint64
	Aux_tail       uint64
	Aux_offset     uint64
	Aux_size       uint64
}

const (
	PerfBitDisabled               uint64 = CBitFieldMaskBit0
	PerfBitInherit                       = CBitFieldMaskBit1
	PerfBitPinned                        = CBitFieldMaskBit2
	PerfBitExclusive                     = CBitFieldMaskBit3
	PerfBitExcludeUser                   = CBitFieldMaskBit4
	PerfBitExcludeKernel                 = CBitFieldMaskBit5
	PerfBitExcludeHv                     = CBitFieldMaskBit6
	PerfBitExcludeIdle                   = CBitFieldMaskBit7
	PerfBitMmap                          = CBitFieldMaskBit8
	PerfBitComm                          = CBitFieldMaskBit9
	PerfBitFreq                          = CBitFieldMaskBit10
	PerfBitInheritStat                   = CBitFieldMaskBit11
	PerfBitEnableOnExec                  = CBitFieldMaskBit12
	PerfBitTask                          = CBitFieldMaskBit13
	PerfBitWatermark                     = CBitFieldMaskBit14
	PerfBitPreciseIPBit1                 = CBitFieldMaskBit15
	PerfBitPreciseIPBit2                 = CBitFieldMaskBit16
	PerfBitMmapData                      = CBitFieldMaskBit17
	PerfBitSampleIDAll                   = CBitFieldMaskBit18
	PerfBitExcludeHost                   = CBitFieldMaskBit19
	PerfBitExcludeGuest                  = CBitFieldMaskBit20
	PerfBitExcludeCallchainKernel        = CBitFieldMaskBit21
	PerfBitExcludeCallchainUser          = CBitFieldMaskBit22
	PerfBitMmap2                         = CBitFieldMaskBit23
	PerfBitCommExec                      = CBitFieldMaskBit24
	PerfBitUseClockID                    = CBitFieldMaskBit25
	PerfBitContextSwitch                 = CBitFieldMaskBit26
)

const (
	PERF_TYPE_HARDWARE   = 0x0
	PERF_TYPE_SOFTWARE   = 0x1
	PERF_TYPE_TRACEPOINT = 0x2
	PERF_TYPE_HW_CACHE   = 0x3
	PERF_TYPE_RAW        = 0x4
	PERF_TYPE_BREAKPOINT = 0x5

	PERF_COUNT_HW_CPU_CYCLES              = 0x0
	PERF_COUNT_HW_INSTRUCTIONS            = 0x1
	PERF_COUNT_HW_CACHE_REFERENCES        = 0x2
	PERF_COUNT_HW_CACHE_MISSES            = 0x3
	PERF_COUNT_HW_BRANCH_INSTRUCTIONS     = 0x4
	PERF_COUNT_HW_BRANCH_MISSES           = 0x5
	PERF_COUNT_HW_BUS_CYCLES              = 0x6
	PERF_COUNT_HW_STALLED_CYCLES_FRONTEND = 0x7
	PERF_COUNT_HW_STALLED_CYCLES_BACKEND  = 0x8
	PERF_COUNT_HW_REF_CPU_CYCLES          = 0x9

	PERF_COUNT_HW_CACHE_L1D  = 0x0
	PERF_COUNT_HW_CACHE_L1I  = 0x1
	PERF_COUNT_HW_CACHE_LL   = 0x2
	PERF_COUNT_HW_CACHE_DTLB = 0x3
	PERF_COUNT_HW_CACHE_ITLB = 0x4
	PERF_COUNT_HW_CACHE_BPU  = 0x5
	PERF_COUNT_HW_CACHE_NODE = 0x6

	PERF_COUNT_HW_CACHE_OP_READ     = 0x0
	PERF_COUNT_HW_CACHE_OP_WRITE    = 0x1
	PERF_COUNT_HW_CACHE_OP_PREFETCH = 0x2

	PERF_COUNT_HW_CACHE_RESULT_ACCESS = 0x0
	PERF_COUNT_HW_CACHE_RESULT_MISS   = 0x1

	PERF_COUNT_SW_CPU_CLOCK        = 0x0
	PERF_COUNT_SW_TASK_CLOCK       = 0x1
	PERF_COUNT_SW_PAGE_FAULTS      = 0x2
	PERF_COUNT_SW_CONTEXT_SWITCHES = 0x3
	PERF_COUNT_SW_CPU_MIGRATIONS   = 0x4
	PERF_COUNT_SW_PAGE_FAULTS_MIN  = 0x5
	PERF_COUNT_SW_PAGE_FAULTS_MAJ  = 0x6
	PERF_COUNT_SW_ALIGNMENT_FAULTS = 0x7
	PERF_COUNT_SW_EMULATION_FAULTS = 0x8
	PERF_COUNT_SW_DUMMY            = 0x9

	PERF_SAMPLE_IP           = 0x1
	PERF_SAMPLE_TID          = 0x2
	PERF_SAMPLE_TIME         = 0x4
	PERF_SAMPLE_ADDR         = 0x8
	PERF_SAMPLE_READ         = 0x10
	PERF_SAMPLE_CALLCHAIN    = 0x20
	PERF_SAMPLE_ID           = 0x40
	PERF_SAMPLE_CPU          = 0x80
	PERF_SAMPLE_PERIOD       = 0x100
	PERF_SAMPLE_STREAM_ID    = 0x200
	PERF_SAMPLE_RAW          = 0x400
	PERF_SAMPLE_BRANCH_STACK = 0x800

	PERF_SAMPLE_BRANCH_USER       = 0x1
	PERF_SAMPLE_BRANCH_KERNEL     = 0x2
	PERF_SAMPLE_BRANCH_HV         = 0x4
	PERF_SAMPLE_BRANCH_ANY        = 0x8
	PERF_SAMPLE_BRANCH_ANY_CALL   = 0x10
	PERF_SAMPLE_BRANCH_ANY_RETURN = 0x20
	PERF_SAMPLE_BRANCH_IND_CALL   = 0x40

	PERF_FORMAT_TOTAL_TIME_ENABLED = 0x1
	PERF_FORMAT_TOTAL_TIME_RUNNING = 0x2
	PERF_FORMAT_ID                 = 0x4
	PERF_FORMAT_GROUP              = 0x8

	PERF_RECORD_MMAP       = 0x1
	PERF_RECORD_LOST       = 0x2
	PERF_RECORD_COMM       = 0x3
	PERF_RECORD_EXIT       = 0x4
	PERF_RECORD_THROTTLE   = 0x5
	PERF_RECORD_UNTHROTTLE = 0x6
	PERF_RECORD_FORK       = 0x7
	PERF_RECORD_READ       = 0x8
	PERF_RECORD_SAMPLE     = 0x9

	PERF_CONTEXT_HV     = -0x20
	PERF_CONTEXT_KERNEL = -0x80
	PERF_CONTEXT_USER   = -0x200

	PERF_CONTEXT_GUEST        = -0x800
	PERF_CONTEXT_GUEST_KERNEL = -0x880
	PERF_CONTEXT_GUEST_USER   = -0xa00

	PERF_FLAG_FD_NO_GROUP = 0x1
	PERF_FLAG_FD_OUTPUT   = 0x2
	PERF_FLAG_PID_CGROUP  = 0x4
)

const (
	CBitFieldMaskBit0  = 0x1
	CBitFieldMaskBit1  = 0x2
	CBitFieldMaskBit2  = 0x4
	CBitFieldMaskBit3  = 0x8
	CBitFieldMaskBit4  = 0x10
	CBitFieldMaskBit5  = 0x20
	CBitFieldMaskBit6  = 0x40
	CBitFieldMaskBit7  = 0x80
	CBitFieldMaskBit8  = 0x100
	CBitFieldMaskBit9  = 0x200
	CBitFieldMaskBit10 = 0x400
	CBitFieldMaskBit11 = 0x800
	CBitFieldMaskBit12 = 0x1000
	CBitFieldMaskBit13 = 0x2000
	CBitFieldMaskBit14 = 0x4000
	CBitFieldMaskBit15 = 0x8000
	CBitFieldMaskBit16 = 0x10000
	CBitFieldMaskBit17 = 0x20000
	CBitFieldMaskBit18 = 0x40000
	CBitFieldMaskBit19 = 0x80000
	CBitFieldMaskBit20 = 0x100000
	CBitFieldMaskBit21 = 0x200000
	CBitFieldMaskBit22 = 0x400000
	CBitFieldMaskBit23 = 0x800000
	CBitFieldMaskBit24 = 0x1000000
	CBitFieldMaskBit25 = 0x2000000
	CBitFieldMaskBit26 = 0x4000000
	CBitFieldMaskBit27 = 0x8000000
	CBitFieldMaskBit28 = 0x10000000
	CBitFieldMaskBit29 = 0x20000000
	CBitFieldMaskBit30 = 0x40000000
	CBitFieldMaskBit31 = 0x80000000
	CBitFieldMaskBit32 = 0x100000000
	CBitFieldMaskBit33 = 0x200000000
	CBitFieldMaskBit34 = 0x400000000
	CBitFieldMaskBit35 = 0x800000000
	CBitFieldMaskBit36 = 0x1000000000
	CBitFieldMaskBit37 = 0x2000000000
	CBitFieldMaskBit38 = 0x4000000000
	CBitFieldMaskBit39 = 0x8000000000
	CBitFieldMaskBit40 = 0x10000000000
	CBitFieldMaskBit41 = 0x20000000000
	CBitFieldMaskBit42 = 0x40000000000
	CBitFieldMaskBit43 = 0x80000000000
	CBitFieldMaskBit44 = 0x100000000000
	CBitFieldMaskBit45 = 0x200000000000
	CBitFieldMaskBit46 = 0x400000000000
	CBitFieldMaskBit47 = 0x800000000000
	CBitFieldMaskBit48 = 0x1000000000000
	CBitFieldMaskBit49 = 0x2000000000000
	CBitFieldMaskBit50 = 0x4000000000000
	CBitFieldMaskBit51 = 0x8000000000000
	CBitFieldMaskBit52 = 0x10000000000000
	CBitFieldMaskBit53 = 0x20000000000000
	CBitFieldMaskBit54 = 0x40000000000000
	CBitFieldMaskBit55 = 0x80000000000000
	CBitFieldMaskBit56 = 0x100000000000000
	CBitFieldMaskBit57 = 0x200000000000000
	CBitFieldMaskBit58 = 0x400000000000000
	CBitFieldMaskBit59 = 0x800000000000000
	CBitFieldMaskBit60 = 0x1000000000000000
	CBitFieldMaskBit61 = 0x2000000000000000
	CBitFieldMaskBit62 = 0x4000000000000000
	CBitFieldMaskBit63 = 0x8000000000000000
)

type SockaddrStorage struct {
	Family uint16
	_      [118]uint8
	_      uint64
}

type TCPMD5Sig struct {
	Addr      SockaddrStorage
	Flags     uint8
	Prefixlen uint8
	Keylen    uint16
	_         uint32
	Key       [80]uint8
}

type HDDriveCmdHdr struct {
	Command uint8
	Number  uint8
	Feature uint8
	Count   uint8
}

type HDGeometry struct {
	Heads     uint8
	Sectors   uint8
	Cylinders uint16
	_         [4]byte
	Start     uint64
}

type HDDriveID struct {
	Config         uint16
	Cyls           uint16
	Reserved2      uint16
	Heads          uint16
	Track_bytes    uint16
	Sector_bytes   uint16
	Sectors        uint16
	Vendor0        uint16
	Vendor1        uint16
	Vendor2        uint16
	Serial_no      [20]uint8
	Buf_type       uint16
	Buf_size       uint16
	Ecc_bytes      uint16
	Fw_rev         [8]uint8
	Model          [40]uint8
	Max_multsect   uint8
	Vendor3        uint8
	Dword_io       uint16
	Vendor4        uint8
	Capability     uint8
	Reserved50     uint16
	Vendor5        uint8
	TPIO           uint8
	Vendor6        uint8
	TDMA           uint8
	Field_valid    uint16
	Cur_cyls       uint16
	Cur_heads      uint16
	Cur_sectors    uint16
	Cur_capacity0  uint16
	Cur_capacity1  uint16
	Multsect       uint8
	Multsect_valid uint8
	Lba_capacity   uint32
	Dma_1word      uint16
	Dma_mword      uint16
	Eide_pio_modes uint16
	Eide_dma_min   uint16
	Eide_dma_time  uint16
	Eide_pio       uint16
	Eide_pio_iordy uint16
	Words69_70     [2]uint16
	Words71_74     [4]uint16
	Queue_depth    uint16
	Words76_79     [4]uint16
	Major_rev_num  uint16
	Minor_rev_num  uint16
	Command_set_1  uint16
	Command_set_2  uint16
	Cfsse          uint16
	Cfs_enable_1   uint16
	Cfs_enable_2   uint16
	Csf_default    uint16
	Dma_ultra      uint16
	Trseuc         uint16
	TrsEuc         uint16
	CurAPMvalues   uint16
	Mprc           uint16
	Hw_config      uint16
	Acoustic       uint16
	Msrqs          uint16
	Sxfert         uint16
	Sal            uint16
	Spg            uint32
	Lba_capacity_2 uint64
	Words104_125   [22]uint16
	Last_lun       uint16
	Word127        uint16
	Dlf            uint16
	Csfo           uint16
	Words130_155   [26]uint16
	Word156        uint16
	Words157_159   [3]uint16
	Cfa_power      uint16
	Words161_175   [15]uint16
	Words176_205   [30]uint16
	Words206_254   [49]uint16
	Integrity_word uint16
}

type Statfs_t struct {
	Type    int64
	Bsize   int64
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Fsid    Fsid
	Namelen int64
	Frsize  int64
	Flags   int64
	Spare   [4]int64
}

const (
	ST_MANDLOCK    = 0x40
	ST_NOATIME     = 0x400
	ST_NODEV       = 0x4
	ST_NODIRATIME  = 0x800
	ST_NOEXEC      = 0x8
	ST_NOSUID      = 0x2
	ST_RDONLY      = 0x1
	ST_RELATIME    = 0x1000
	ST_SYNCHRONOUS = 0x10
)

type TpacketHdr struct {
	Status  uint64
	Len     uint32
	Snaplen uint32
	Mac     uint16
	Net     uint16
	Sec     uint32
	Usec    uint32
	_       [4]byte
}

type Tpacket2Hdr struct {
	Status    uint32
	Len       uint32
	Snaplen   uint32
	Mac       uint16
	Net       uint16
	Sec       uint32
	Nsec      uint32
	Vlan_tci  uint16
	Vlan_tpid uint16
	_         [4]uint8
}

type Tpacket3Hdr struct {
	Next_offset uint32
	Sec         uint32
	Nsec        uint32
	Snaplen     uint32
	Len         uint32
	Status      uint32
	Mac         uint16
	Net         uint16
	Hv1         TpacketHdrVariant1
	_           [8]uint8
}

type TpacketHdrVariant1 struct {
	Rxhash    uint32
	Vlan_tci  uint32
	Vlan_tpid uint16
	_         uint16
}

type TpacketBlockDesc struct {
	Version uint32
	To_priv uint32
	Hdr     [40]byte
}

type TpacketReq struct {
	Block_size uint32
	Block_nr   uint32
	Frame_size uint32
	Frame_nr   uint32
}

type TpacketReq3 struct {
	Block_size       uint32
	Block_nr         uint32
	Frame_size       uint32
	Frame_nr         uint32
	Retire_blk_tov   uint32
	Sizeof_priv      uint32
	Feature_req_word uint32
}

type TpacketStats struct {
	Packets uint32
	Drops   uint32
}

type TpacketStatsV3 struct {
	Packets      uint32
	Drops        uint32
	Freeze_q_cnt uint32
}

type TpacketAuxdata struct {
	Status    uint32
	Len       uint32
	Snaplen   uint32
	Mac       uint16
	Net       uint16
	Vlan_tci  uint16
	Vlan_tpid uint16
}

const (
	TPACKET_V1 = 0x0
	TPACKET_V2 = 0x1
	TPACKET_V3 = 0x2
)

const (
	SizeofTpacketHdr  = 0x20
	SizeofTpacket2Hdr = 0x20
	SizeofTpacket3Hdr = 0x30
)
