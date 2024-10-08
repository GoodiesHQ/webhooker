# webhooker
webhook -> webhook(s)

This repository is an incredibly simple webhook proxy and forwarder. It addresses two issues I ran into using various clients that support webhook notifications:


1) What if I want to take one event and send it to multiple webhooks at once without implementing multiple notification profiles?
2) What if the webhook URL provided by the target contains characters the client can't process (?, %, etc...)?

Enter WebHooker... the configuration is simple. Set a name (endpoint URL, case sensitive) and set your targets.

```yaml
webhooks:
  - name: abcdefg
    targets:
      - https://tenant1.microsoft.com/teams/channel/1234567?apiver=1.2.3&sig=aabbccddeeff00
      - https://tenant2.microsoft.com/teams/channel/7654321?apiver=1.2.3&sig=a1b2c3d4f5g677
```

This means that sending a webhook to the URL:
```
webhook.com/abcdef
```
it will receive all of the post and query data, rebuild the request, and send it to all of the upstream targets simultaneously.


#### webhooker is not smart
It just forwards the requests. The same query parameters, the same post data, the same request headers. Only in the logs can you determine if the webhook was successful to the targets. The client side currently has no visibility to that. Otherwise, it will simply always respond with 200 OK no matter what.

## Responses:

| Status Code | Meaning | 
|---------|---------------------------------------------------------------------|
| 200 OK | webhooker forwarded the requests (no indication of remote statuses)|
| 404 Not Found | webhooker has no endpoint configured for that URL |
| 405 Not Allowed | webhooker received something other than a POST request |
