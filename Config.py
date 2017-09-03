import os
import logging
import yaml


class Config:

    @classmethod
    def load(self):
        config_dir = os.path.join(
            os.path.dirname(os.path.realpath(__file__)),
            'config'
        )

        config = {}

        for config_name in ['default', 'local']:
            fname = 'config.{}.yml'.format(config_name)
            path = os.path.join(config_dir, fname)
            if os.path.exists(path):
                data = load_config(path)
                logging.debug('loading config at path: %s', path)
                merge_rec(src=data, dst=config)

        return config


def load_config(fname):
    """
    load config at fname, assuming it's a yml file
    """
    with open(fname, 'r') as f:
        data = yaml.load(f)
    return data


def merge_rec(src, dst):
    """
    merge recursively when both src and dst are dictionaries.
    otherwise, src overrides dst
    """
    if type(src) == dict and type(dst) == dict:
        for k, src_v in src.items():
            dst_v = dst.get(k)
            if dst_v is not None:
                v = merge_rec(src=src_v, dst=dst_v)
            else:
                v = src_v
            dst[k] = v
        return dst
    else:
        return src
