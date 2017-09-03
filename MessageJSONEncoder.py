from json import JSONEncoder

from envisalink import Message
from Enum import EnumItem


class MessageJSONEncoder(JSONEncoder):

    def default(self, obj):
        if isinstance(obj, Message):
            return self._to_json_message(obj)
        elif isinstance(obj, EnumItem):
            return self._to_json_enum(obj)
        elif isinstance(obj, bytes):
            return obj.decode('ascii')
        else:
            return JSONEncoder.default(self, obj)

    def _to_json_message(self, obj):
        return {
            'uid': obj.uid,
            'code': obj.code,
            'data': obj.data,
        }

    def _to_json_enum(self, obj):
        return obj.name
