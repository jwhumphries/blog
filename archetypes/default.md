---
date: '{{ .Date }}'
draft: true
author: 'John Humphries'
title: '{{ replace .File.ContentBaseName "-" " " | title }}'
description: ''
tags: []
categories: ['{{ .Parent }}']
ShowToc: true
TocOpen: true
---
