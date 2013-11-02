#! /usr/bin/env python
#
# Simple test module for Envisalink 2DS network integration.
# Author: matt@maddogsw.com
#
# Bugs:
# - simple, incomplete
# - should reconnect automatically if connection drops

import optparse
import socket
import sys
import time


class Envisalink:
	def __init__(self):
		self.verbose = False
	
	def connect(self, host_address):
		self.host = host_address
		self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
		self.socket.connect((host_address, 4025))
		self.socket.setblocking(0)

	def set_verbose(self, verbosity):
		self.verbose = verbosity

	def login(self, password):
		self.send_command(005, password)

	def send_command(self, command, data_bytes = []):
		# cmd_bytes = '%03s' % int(command) # BUG: why doesn't this work?
		cmd_bytes = str(command).zfill(3)
		cmd = []
		checksum = 0
		for byte in cmd_bytes:
			cmd.append(byte)
			checksum += ord(byte)
		for byte in data_bytes:
			cmd.append(byte)
			checksum += ord(byte)

		checksum = checksum % 256
		cmd.extend([hex(nibble)[-1].upper() for nibble in [ checksum / 16, checksum % 16]])
		cmd.extend((chr(0x0D), chr(0x0A)))

		if self.verbose:
			print "debug: send command: " + str(cmd)
		
		self.socket.send(''.join(cmd))

	def read_reply(self):
		reply = ''
		try:
			while True: # loop until command ends with CRLF
				chunk = self.socket.recv(1024)
				if self.verbose:
					bytes = [byte for byte in chunk]
					print "debug: receive reply, %d bytes: %s" % (len(chunk), bytes)
				reply = reply + chunk
				if reply[-2:] == '\x0D\x0A':
					break
				# BUG: what happens if we get CRLF in middle of reply... we should find it, and stash the rest for the next invocation
				# Should not depend on network messages arriving in neat-sized chunks! And we need to avoid tripping over CRLF sequences
				# embedded in data or checksum fields, probably by decoding and verifying checksum.
		except socket.error:
			return None
		return reply

if __name__ == '__main__':
		# When invoked directly, parse the command line to find hostname and password,
		# then connect, log in, and enter a monitoring loop which just prints system events
		# until killed.
		
		# parse command line
		p = optparse.OptionParser()
		p.add_option('-H', '--host', default = 'envisalink', help = 'network hostname for envisalink module, default envisalink')
		p.add_option('-P', '--password', default = 'user', help = 'password for envisalink user account, default user')
		p.add_option('-V', '--verbose', action = 'store_true', help = 'enable debug output')
		(options, args) = p.parse_args()
		
		# connect and log in
		e = Envisalink()
		if options.verbose:
			e.set_verbose(True)
		e.connect(options.host)
		e.login(options.password)
		
		# monitor loop
		sleep = 0
		while(True):
			reply = e.read_reply()
			if reply == None:
				sleep += 1
				if sleep == 10:
					e.send_command(1)
					sleep = 0
				else:
					time.sleep(1)
			else:
				print "event: " + reply.rstrip()

