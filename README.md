# Description
This is a project that implement the CRUSH algorithm.
It's mainly doing two things:

1. Construct a physical cluster into logical *tree* structure
2. Selecting specific number of items follows the *rule*

# How to construct your cluster
ceph-crush-ana accept a json of tree representation to help itself construct the cluster.
There are two types of item in structer.
Bucket and osd.

For bucket type, you can represent its structure ,like:
```
{
    "add_ops": [
        {
            "new_items": ["IDC-01"],
            "type": "root"
        },
        {
            "new_items": ["room1", "room2", "room3"],
            "type": "room"
        },
        {
            "new_items": ["vrack1", "vrack2"],
            "type": "rack"
        },
        {
            "new_items": ["hostA", "hostB"],
            "type": "host"
        }
    ],
    "move_ops": [
        {
            "source_names": ["room1", "room2", "room3"],
            "target_name": "EBS_SHB",
            "target_type": "root"
        },
        {
            "source_names": ["vrack1"],
            "target_name": "room1",
            "target_type": "room"
        },
        {
            "source_names": ["vrack2"],
            "target_name": "room2",
            "target_type": "room"
        },
        {
            "source_names": ["hostA"],
            "target_name": "vrack1",
            "target_type": "rack"
        },
        {
            "source_names": ["hostB"],
            "target_name": "vrack2",
            "target_type": "rack"
        }
    ]
}
```

After that, you will get a tree like:

IDC-01 && (root)<br>
&&|-room1 && (room)<br>
&&&& |-vrack1 && (rack) <br>
&&&&&&&& |-hostA && (host) <br>
&&|-room2 && (room)<br>
&&&& |-vrack2 && (rack) <br>
&&&&&&&& |-hostB && (host) <br>
&&|-room3 && (room) <br>
             

For osds, you must represent them under a certain *Host*, like:
```
{
    "osd_num": 3,
    "target_host": "hostA"
}
```
After that, hostA will contain 3 osds, and the tree will become:

IDC-01 && (root)<br>
&&|-room1 && (room)<br>
&&&& |-vrack1 && (rack) <br>
&&&&&&&& |-hostA && (host) <br>
& & & & & & & & |-osd.0 && (osd) <br>
& & & & & & & & |-osd.1 && (osd) <br>
& & & & & & & & |-osd.2 && (osd) <br>
&&|-room2 && (room)<br>
&&&& |-vrack2 && (rack) <br>
&&&&&&&& |-hostB && (host) <br>
&&|-room3 && (room) <br>

The osd name will create by ceph-crush-ana.

# How to represent your Rule

A rule will contain several steps, which finally lead to couples of items you want.

For example:
```
{
    "steps": [
        {"num": 1, "type": "root"},
        {"num": 1, "type": "room"},
        {"num": 3,"type": "host"},
        {"num": 1,"type": "osd"}
    ]
}
```
Then ceph-crush-ana will do following things:
step1: choose the root <br>
step2: select 1 room from root <br>
step3: select 3 hosts from each room of step2 <br>
step4: select 1 osd from each host of step3 <br>

