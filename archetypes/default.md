---
date: '{{ .Date }}'
draft: true
author: 'John Humphries'
title: '{{ replace .File.ContentBaseName "-" " " | title }}'
description: ''
tags: []
categories: ['{{ replace .Dir "/" "" }}']
ShowToc: true
TocOpen: true
---
