#!/bin/bash
set +ex

ansible-playbook -i ../inventory.yaml ../playbook.yaml "$@"