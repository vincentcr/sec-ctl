
class Enum(object):

    def __init__(self, *args, **kwargs):
        
        names = {}
        vals = {}

        for val,name in enumerate(args):
            names[name] = val
            vals[val] = name

        for name,val in kwargs.iteritems():
            names[name] = val
            vals[val] = name            

        for name in names:
            assert(isinstance(name, str) or isinstance(name,unicode))

        for val in vals:
            assert(isinstance(val, int))

        sorted_vals = sorted(vals.items(), key = lambda tpl : tpl[0])

        self._names = []
        self._items = { } 
        self._by_val = { }
        for val,name in sorted_vals:
            self._names.append(name)
            item = EnumItem(name = name, val = val)
            self._items[name] = item
            self._by_val[val] = item 


    def by_val(self, val):
        return self._by_val.get(val)

    def items(self):
        return dict(self._items)


    def __getattr__(self, name):
        if name in self._items:
            return self._items[name]

        raise AttributeError('%s not in %s' % (name, self._names))



class EnumItem(object):
    
    def __init__(self, name, val):
        self.name = name
        self.val = val

    def __repr__(self):
        return self.name
