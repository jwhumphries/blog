---
date: '{{ .Date }}'
draft: true
author: 'John Humphries'
title: '{{ replace .File.ContentBaseName "-" " " | title }}'
description: ''
topics: []
subjects: ['{{ replace .Dir "posts/" "" }}']
---
