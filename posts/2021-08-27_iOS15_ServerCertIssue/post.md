# iOS 15 / XCode 13 Error - Certificate is not permitted for this usage

## Summary

We had an iOS application that was erroring when connecting to its backend API.  Long story short, the application 
uses TLS client authentication to authenticate with the API against an internal PKI CA certificate, so the application is manually calling [SecTrustEvaluateWithError](https://developer.apple.com/documentation/security/2980705-sectrustevaluatewitherror) to verify both the 
server certificate and the chain (ie, that it was signed by the internal CA). 

## Issue

When running the application in the XCode 13 beta or any device with iOS 15 (including the simulator), no API call to the HTTP API would succeed. 

The following code was being run inside the application:
```swift
SecTrustSetAnchorCertificates(serverTrust, [rootCertificateAuthorityCert] as CFArray) //allow our CA
SecTrustSetAnchorCertificatesOnly(serverTrust, false) // also allow regular CAs.

var error:CFError?
let trustResult = SecTrustEvaluateWithError(serverTrust, &error)

if error != nil {
    print(error!.localizedDescription)
}
```
This code verifies the certificate is valid, AND that it's signed by our CA.

The ```trustResult``` would always return false, and the error that was printed out was: **Certificate is not permitted for this usage**

The certificate in question had the following extensions for the CSR: 
```
[ req_ext ]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
```

For some reason, this certificate had two keyUsages defined, *keyEncipherment* and *dataEncipherment*. 
Why? No idea.  Probably a copy/paste from the good ol' internet. 

## The solution

Regenerate the certificate and **remove the keyUsage parameter in the CSR config**. Have nothing, only the extendedKeyUsage.  
Once this was deployed to the serer, the application immediately worked without any issues.  
It appears that Apple has become very specific in iOS 15 of what a server certificate should look like.  Even though the *extendedKeyUsage* was set to ```serverAuth```, it was complaining about the keyUsage that was specified on the certificate. 


That was a fun few hours of debugging.  