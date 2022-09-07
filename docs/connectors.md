---
title: Connectors
nav_order: 3
has_children: true
---

#  Connectors

## Introduction

The Connectors module provides a large number of out-of-the-box, image-based connectors. Source-Connectors are responsible for obtaining data from the source, converting these data in different formats into a unified CloudEvent, and sending them to the target. Sink-Connectors are responsible for performing general processing on the data in CloudEvents format (such as: Adding data to a MySQL database). Connectors enable their users to no longer focus on how to obtain or deal with data formats so they can focus on their business without any worries.

## Source & Sink

As an essential part of Vanus, the most significant role of Connectors are to lower the threshold for users to use Vanus with a rich ecosystem. The right Source-Connector can import any data generated by other platforms and convert it into a standard CloudEvent.
In a similar form, the right Sink-Connector can receive and handle the Cloud Event to accomplish its own distinct action without limitation. [List of Sink and Source Connectors](/list.md)

## vanus & vance Advantage

There is a multitude of companies on the market (Confluent, Airbyte) that have their own systems. When discussing Connectors, everyone will undoubtedly consider this question: The company Confluent already has more than 100 Connectors. Compared with such a rich ecosystem, what advantages or differences can Vanus Connectors offer?

The answer to this question is clear, native support for the CloudEvent Specifications. 
After all, we chose CloudEvent as the data format, it is natural that we have decided to take it as the format for data transmission, because in fact, CE is the norm for events in the cloud-native era. We know and believe it will become a Standard in the cloud-native era, its existence in itself is our differentiating advantage. We have a clear advantage compared to other companies like Confluent, native support for CloudEvents, and code-less benefits. It is a feature that Confluent Connectors or others cannot replicate without a complete restructure of their protocol. 


