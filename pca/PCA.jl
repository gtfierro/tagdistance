f = open("out.csv","w+")
filename = ARGS[1]
println("reading data from ",filename)
data = readdlm(filename,',',String)
datasize = size(data)
tags = Dict{ASCIIString, Uint8}()
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

println("creating data matrix X")
#X = spzeros(size(data,1),length(tags))
X = SharedArray(Uint8, size(data,1), length(tags))
o = convert(Uint8,1)
for i=1:size(data,1)
    mytags = split(data[i,2])
    for tag=mytags
        j = wordindexes[tag]
        X[i,j] = o
    end
end
X = sparse(sdata(X))
m = mean(X,1)
println("mean:",size(m))
U = broadcast(-, X, m)
println("U:",size(U))
cov = *(U',U) / size(U,1)
println("cov:",size(cov))
F = eigfact(cov)
# F's eigenvalues are sorted in ascending order, so we take the last 3
sumeigvals = sum(real(F[:values]))
println("sum all eigenvalues:",sumeigvals)
top3values = F[:values][end-2:end]
println("sum top 3 eigenvalues:", sum(top3values))
top3vectors = F[:vectors][end-2:end, :]'
println("top3eigvec:",size(top3vectors))
projected = *(U,top3vectors)
println(size(projected))

writecsv(f, projected)
close(f)
