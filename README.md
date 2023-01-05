# Testaservice - Convenient HTTP Service assertions based on testify library

This test library provides an HTTP service useful for use in testing allowing full
exercise of the HTTP transport in unit tests.


## Example

```
func TestMyHTTPClient(t *testing.T){
    t.Run("some sort of client test", func(t *testing.T){
        ts := NewTestService(t)
        // Do some client things

        var request ClientRequest
        ts.AssertCalled()
        ts.AssertReceivedAs(&request)
        // etc
    })
}
```
