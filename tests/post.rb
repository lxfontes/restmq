require "net/http"
require "uri"
require 'json'

uri = URI.parse("http://localhost:8080/queue")

cmd = {
    :cmd => "add",
    :queue => "tque",
    :value => "pizza"
}
response = Net::HTTP.post_form(uri, {"body" => cmd.to_json})
puts(response.body)
