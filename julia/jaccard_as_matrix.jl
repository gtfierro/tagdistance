filename = ARGS[1]
println("using file ",filename)
println("reading data from ",filename)
data = readdlm(filename,',',String)
datasize = size(data)
println("finished reading data ")


tags = Dict{ASCIIString, Int64}()
for i=1:size(data,1)
    words = split(data[i,2])
    for word in words
        if haskey(tags, word)
            tags[word] += 1
        else
            tags[word] = 1
        end
    end
end
words = [x for x in keys(tags)]
wordindexes = [words[i] => i for i=1:length(words)]

function vector(i)
    words = Set(split(data[i,2])...)
    [x in words ? 1 : 0 for x in keys(tags)]
end

println("creating data matrix X")
X = spzeros(size(data,1),length(tags))
for i=1:size(data,1)
    mytags = split(data[i,2])
    for tag=mytags
        j = wordindexes[tag]
        X[i,j] = 1
    end
end

Xt = sparse(X')
ONE = ones(size(X))
println("computing intersection")
intersection = *(X,Xt)
Xnonzero = sparse(max(intersection ./ intersection,0))
println("computing union")
union = 2 * (*(X,ONE')) - intersection
println("computing jaccard")
jaccard = sparse(intersection ./ union)
println(nfilled(jaccard))
println(size(jaccard))
println(data[1,1],data[1,2])
println(data[50,1],data[50,2])
println(intersection[1,50])
println(union[1,50])
println(jaccard[1,50])
