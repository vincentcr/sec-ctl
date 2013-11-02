#!/usr/bin/env python 

import os, sys, re

def main():
    #print encode_cmd(654, 3)
    tpi = Tpi()

    tpi.send_cmd(5, 'Q4m1gh')
    tpi.send_cmd(1)

    return


class Tpi(object):

    def __init__(self):
        self.socket = sys.stdout


    def send_cmd(self, code, data = None):
        cmd = self.encode_cmd(code, data)
        self.socket.write(cmd)

    def encode_cmd(self, code, data):
        code_str = '%03d' % code
        if data is not None:
            data_str = str(data)
        else:
            data_str = ''
        cmd = code_str + data_str
        cmd_chars = [ ord(c) for c in cmd ]
        checksum = sum(cmd_chars) & 0xff
        checksum_hex = '%x' % checksum
        encoded = code_str + data_str + checksum_hex + '\n\r'
        return encoded




main()



