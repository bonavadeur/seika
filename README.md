# Seika - 星歌

[![LICENSE](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

Kubernetes Custom Resource maintains quantity of Pods in each Node

![](images/seika.jpg)

## 1. Introduction

Seika is a Kubernetes Custom Resource maintains quantity of Pods in each Node without creating many ReplicaSet or Deployment.

This document is for **use purpose** only. All you need is in `dist/`

## 2. Install

Install CRD to your Kubernetes System

```bash 
$ kubectl apply -f dist/install.yaml
```

## 3. How to use

Take a quick view in `dist/seika-sample.yaml`

```yaml
apiVersion: batch.bonavadeur.io/v1
kind: Seika
metadata:
  labels:
    app.kubernetes.io/name: seika
    app.kubernetes.io/instance: seika-sample
    app.kubernetes.io/part-of: seika
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: seika
  name: seika-sample
spec:
  repurika:
    node1: 1
    node2: 2
    node3: 3
  selector:
    matchLabels:
      bonavadeur.io/seika: seika-sample
  template:
    spec:
      containers:
      - name: shuka
        image: docker.io/bonavadeur/shuka:sleep
```

Note: 
+ Pod's labels is not required. Pods are labeled with Seika's Labels. In the yaml file above, Pods are labels with `bonavadeur.io/seika: seika-sample`  
+ `.spec.repurika` field specify quantity of Pods in each node, with key is nodename and value is an integer specify desired number of Pods in that node

Let's apply a Seika Instance

```bash
$ kubectl apply -f dist/seika-sample.yaml
$ kubectl get pod | grep seika
seika-sample-node1-u8knc       1/1     Running   0               15m    10.233.102.166   node1   <none>           <none>
seika-sample-node2-cgo5g       1/1     Running   0               16m    10.233.75.25     node2   <none>           <none>
seika-sample-node2-qkp5u       1/1     Running   0               16m    10.233.75.47     node2   <none>           <none>
seika-sample-node3-7qf9g       1/1     Running   0               8s     10.233.71.45     node3   <none>           <none>
seika-sample-node3-efl0y       1/1     Running   0               9s     10.233.71.49     node3   <none>           <none>
seika-sample-node3-q8i96       1/1     Running   0               16m    10.233.71.43     node3   <none>           <none>
```

## 4. Contributeur

Đào Hiệp - Bonavadeur - ボナちゃん  
The Future Internet Laboratory, Room E711 C7 Building, Hanoi University of Science and Technology, Vietnam.
未来のインターネット研究室, C7 の E ７１１、ハノイ百科大学、ベトナム。  

![](images/github-wp.png)
