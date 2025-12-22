# Set up k8s cluster on Cloudlab

* Change the `config.json` to include all hostnames and username
* Run `python3 setup_kube.py` which will install everything and set up a k8s cluster. The first node in the hostname list will be the control node, and the rest will be worker nodes.